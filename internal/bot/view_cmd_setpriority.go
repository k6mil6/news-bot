package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k6mil6/news-bot/internal/botkit"
	"github.com/k6mil6/news-bot/internal/state"
)

func ViewCmdSetPriority(stateMachine *state.Machine) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		var reply = tgbotapi.NewMessage(update.Message.Chat.ID, "Укажите ID источника и его приоритет (через пробел)")

		stateMachine.Set(update.Message.Chat.ID, state.WaitingForSourceIDAndPriority)

		_, err := bot.Send(reply)

		return err
	}
}

//func ViewCmdSetPriority(prioritySetter PrioritySetter) botkit.ViewFunc {
//	type setPriorityArgs struct {
//		SourceID int64 `json:"source_id"`
//		Priority int   `json:"priority"`
//	}
//	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
//		args, err := botkit.ParseJSON[setPriorityArgs](update.Message.CommandArguments())
//		if err != nil {
//			return err
//		}
//
//		if err := prioritySetter.SetPriority(ctx, args.SourceID, args.Priority); err != nil {
//			return err
//		}
//
//		var reply = tgbotapi.NewMessage(update.Message.Chat.ID, "Приоритет установлен")
//		reply.ParseMode = "MarkdownV2"
//
//		if _, err := bot.Send(reply); err != nil {
//			return err
//		}
//
//		return nil
//	}
//}
