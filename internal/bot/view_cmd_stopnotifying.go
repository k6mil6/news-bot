package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k6mil6/news-bot/internal/botkit"
)

type Stopper interface {
	Pause() error
}

func ViewCmdStopNotifying(stopper Stopper) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if err := stopper.Pause(); err != nil {
			msg := tgbotapi.NewMessage(update.FromChat().ID, "Постинг статей уже приостановлен!")
			_, err := bot.Send(msg)
			return err
		}
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Постинг статей приостановлен")
		_, err := bot.Send(msg)
		return err
	}
}
