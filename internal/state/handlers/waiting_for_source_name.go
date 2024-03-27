package handlers

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k6mil6/news-bot/internal/botkit"
	"github.com/k6mil6/news-bot/internal/model"
	"github.com/k6mil6/news-bot/internal/state"
)

type SourceSaver interface {
	Add(ctx context.Context, source model.Source) (int64, error)
}

func ViewWaitingForSourceName(sourceSaver SourceSaver, stateMachine *state.Machine) botkit.StateViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		src := model.Source{Name: update.Message.Text}

		sourceID, err := sourceSaver.Add(ctx, src)
		if err != nil {
			return err
		}

		var (
			msgText = fmt.Sprintf(
				"ID источника: `%d`\\. Отправьте URL источника\\.",
				sourceID,
			)
			reply = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		)

		reply.ParseMode = "MarkdownV2"

		stateMachine.Set(update.Message.Chat.ID, state.WaitingForSourceURL)

		_, err = bot.Send(reply)
		return err
	}
}
