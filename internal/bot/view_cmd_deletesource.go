package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k6mil6/news-bot/internal/botkit"
	"strconv"
)

type SourceDeleter interface {
	Delete(context context.Context, id int64) error
}

func ViewCmdDeleteSource(deleter SourceDeleter) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		idArg := update.Message.CommandArguments()

		id, err := strconv.ParseInt(idArg, 10, 64)
		if err != nil {
			if _, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат ID")); err != nil {
				return err
			}
			return err
		}

		if err = deleter.Delete(ctx, id); err != nil {
			if _, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Источник не найден")); err != nil {
				return err
			}
			return err
		}

		if _, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Источник удален")); err != nil {
			return err
		}

		return nil
	}
}
