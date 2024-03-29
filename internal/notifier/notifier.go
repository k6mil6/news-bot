package notifier

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-shiori/go-readability"
	"github.com/k6mil6/news-bot/internal/botkit/markup"
	"github.com/k6mil6/news-bot/internal/model"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ArticleProvider interface {
	AllNotPosted(ctx context.Context, limit uint64) ([]model.Article, error)
	MarkAsPosted(ctx context.Context, article model.Article) error
}

type Summariser interface {
	Summarise(text string) (string, error)
}

type Notifier struct {
	articles         ArticleProvider
	summariser       Summariser
	bot              *tgbotapi.BotAPI
	sendInterval     time.Duration
	lookupTimeWindow time.Duration
	channelId        int64
	pauseCond        *sync.Cond
	paused           bool
}

func New(
	articleProvider ArticleProvider,
	summariser Summariser,
	bot *tgbotapi.BotAPI,
	sendInterval time.Duration,
	lookupTimeWindow time.Duration,
	channelId int64,
) *Notifier {
	return &Notifier{
		articles:         articleProvider,
		summariser:       summariser,
		bot:              bot,
		sendInterval:     sendInterval,
		lookupTimeWindow: lookupTimeWindow,
		channelId:        channelId,
		pauseCond:        sync.NewCond(new(sync.Mutex)),
	}
}

func (n *Notifier) Start(ctx context.Context) error {
	log.Println("[INFO] Starting notifier...")
	ticker := time.NewTicker(n.sendInterval)
	defer ticker.Stop()

	if err := n.SelectAndSendArticle(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			n.pauseCond.L.Lock()
			for n.paused {
				n.pauseCond.Wait()
			}
			n.pauseCond.L.Unlock()
			if err := n.SelectAndSendArticle(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (n *Notifier) Pause() error {
	n.pauseCond.L.Lock()
	defer n.pauseCond.L.Unlock()

	if n.paused {
		return errors.New("already paused")
	}

	n.paused = true
	return nil
}

func (n *Notifier) Resume() error {
	n.pauseCond.L.Lock()
	defer n.pauseCond.L.Unlock()

	if !n.paused {
		return errors.New("already running")
	}
	n.paused = false
	n.pauseCond.Signal()

	return nil
}

func (n *Notifier) StopNotifyingFor(duration time.Duration) {
	n.sendInterval = duration
}

func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	topOneArticles, err := n.articles.AllNotPosted(ctx, 1)
	if err != nil {
		return err
	}

	if len(topOneArticles) == 0 {
		return nil
	}

	article := topOneArticles[0]

	summary, err := n.extractSummary(article)
	if err != nil {
		return err
	}

	if err := n.SendArticle(article, summary); err != nil {
		return err
	}

	return n.articles.MarkAsPosted(ctx, article)
}

func (n *Notifier) extractSummary(article model.Article) (string, error) {
	var r io.Reader

	if article.Summary != "" {
		r = strings.NewReader(article.Summary)
	} else {
		resp, err := http.Get(article.Link)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		r = resp.Body
	}

	doc, err := readability.FromReader(r, nil)
	if err != nil {
		return "", err
	}

	summary, err := n.summariser.Summarise(cleanupText(doc.TextContent))
	if err != nil {
		return "", err
	}

	return "\n\n" + summary, nil
}

var redundantNewLines = regexp.MustCompile(`\n{3,}`)

const readMoreText = "\n\nЧитать далее ->"

func cleanupText(text string) string {
	cleanText := redundantNewLines.ReplaceAllString(text, "\n")
	return CutStringToWords(strings.Replace(cleanText, "Читать далее", readMoreText, 1), 200)
}

func CutStringToWords(str string, numWords int) string {
	words := strings.Fields(str)
	if len(words) <= numWords {
		return str + readMoreText
	}
	return strings.Join(words[:numWords], " ")
}

func (n *Notifier) SendArticle(article model.Article, summary string) error {
	const msgFormat = "*%s*%s\n\n%s"

	msg := tgbotapi.NewMessage(n.channelId, fmt.Sprintf(
		msgFormat,
		markup.EscapeForMarkdown(article.Title),
		markup.EscapeForMarkdown(summary),
		markup.EscapeForMarkdown(article.Link),
	))
	msg.ParseMode = "MarkdownV2"

	_, err := n.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
