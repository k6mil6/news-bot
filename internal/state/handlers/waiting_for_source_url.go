package handlers

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k6mil6/news-bot/internal/botkit"
	"github.com/k6mil6/news-bot/internal/state"
)

type SourceURLSetter interface {
	GetLatestID(ctx context.Context) (int64, error)
	SetURL(ctx context.Context, sourceID int64, url string) error
}

func ViewWaitingForSourceURL(sourceSetter SourceURLSetter, stateMachine *state.Machine) botkit.StateViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		sourceURL := update.Message.Text

		sourceID, err := sourceSetter.GetLatestID(ctx)
		if err != nil {
			return err
		}

		if err := sourceSetter.SetURL(ctx, sourceID, sourceURL); err != nil {
			return err
		}

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Отправьте приоритет источника")

		stateMachine.Set(update.Message.Chat.ID, state.WaitingForSourcePriority)

		_, err = bot.Send(reply)
		return err
	}
}
