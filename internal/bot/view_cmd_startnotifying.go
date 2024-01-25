package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k6mil6/news-bot/internal/botkit"
)

type Starter interface {
	Resume() error
}

func ViewCmdStartNotifying(starter Starter) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if err := starter.Resume(); err != nil {
			msg := tgbotapi.NewMessage(update.FromChat().ID, "Постинг статей уже запущен!")
			_, err := bot.Send(msg)
			return err
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, "Постинг статей возобновлен")
		_, err := bot.Send(msg)
		return err
	}
}
