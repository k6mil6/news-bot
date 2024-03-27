package handlers

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k6mil6/news-bot/internal/botkit"
	"strconv"
	"strings"
)

type PrioritySetter interface {
	SetPriority(ctx context.Context, sourceID int64, priority int) error
}

func ViewWaitingForSourceIDAndPriority(prioritySetter PrioritySetter) botkit.StateViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		msgText := strings.Split(update.Message.Text, " ")

		if len(msgText) != 2 {
			return nil
		}

		sourceID, err := strconv.ParseInt(msgText[0], 10, 64)
		if err != nil {
			return err
		}
		priority, err := strconv.Atoi(msgText[1])
		if err != nil {
			return err
		}

		if err := prioritySetter.SetPriority(ctx, sourceID, priority); err != nil {
			return err
		}

		var reply = tgbotapi.NewMessage(update.Message.Chat.ID, "Приоритет установлен")

		_, err = bot.Send(reply)
		return err
	}
}
