package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k6mil6/news-bot/internal/botkit"
	"github.com/k6mil6/news-bot/internal/state"
)

func ViewCmdStart(stateMachine *state.Machine) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Hello!")
		stateMachine.Set(update.FromChat().ID, state.None)
		_, err := bot.Send(msg)
		return err
	}
}
