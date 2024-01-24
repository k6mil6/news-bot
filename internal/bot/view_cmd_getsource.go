package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k6mil6/news-bot/internal/botkit"
	"github.com/k6mil6/news-bot/internal/model"
	"strconv"
)

type SourceProvider interface {
	SourceByID(ctx context.Context, id int64) (*model.Source, error)
}

func ViewCmdGetSource(provider SourceProvider) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		idArg := update.Message.CommandArguments()

		id, err := strconv.ParseInt(idArg, 10, 64)
		if err != nil {
			if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат ID")); err != nil {
				return err
			}
			return err
		}

		source, err := provider.SourceByID(ctx, id)
		if err != nil {
			return err
		}

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, formatSource(*source))
		reply.ParseMode = "MarkdownV2"

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
