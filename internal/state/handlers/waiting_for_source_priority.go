package handlers

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k6mil6/news-bot/internal/botkit"
	"github.com/k6mil6/news-bot/internal/state"
	"strconv"
)

type SourcePrioritySetter interface {
	GetLatestID(ctx context.Context) (int64, error)
	SetPriority(ctx context.Context, sourceID int64, priority int) error
}

func ViewWaitingForSourcePriority(sourceSetter SourcePrioritySetter, stateMachine *state.Machine) botkit.StateViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		priority, err := strconv.Atoi(update.Message.Text)
		if err != nil {
			return err
		}

		sourceID, err := sourceSetter.GetLatestID(ctx)
		if err != nil {
			return err
		}

		if err := sourceSetter.SetPriority(ctx, sourceID, priority); err != nil {
			return err
		}

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Источник добавлен")

		stateMachine.Set(update.Message.Chat.ID, state.None)

		_, err = bot.Send(reply)
		return nil
	}
}
