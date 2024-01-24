package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k6mil6/news-bot/internal/botkit"
	"time"
)

type NotifierStopper interface {
	StopNotifyingFor(duration time.Duration)
}

func ViewCmdStopNotifyingFor(stopper NotifierStopper) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		duration, err := time.ParseDuration(update.Message.CommandArguments())
		if err != nil {
			return err
		}
		stopper.StopNotifyingFor(duration)

		msg := tgbotapi.NewMessage(update.FromChat().ID, fmt.Sprintf("Постинг статей приостановлен на %s", duration))
		_, err = bot.Send(msg)
		return err
	}
}
