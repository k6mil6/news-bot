package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/k6mil6/news-bot/internal/botkit"
	"github.com/k6mil6/news-bot/internal/model"
	"github.com/k6mil6/news-bot/internal/state"
)

type SourceStorage interface {
	Add(ctx context.Context, source model.Source) (int64, error)
}

func ViewCmdAddSource(stateMachine *state.Machine) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		reply := tgbotapi.NewMessage(update.FromChat().ID, "Введите название источника")
		stateMachine.Set(update.FromChat().ID, state.WaitingForSourceName)

		_, err := bot.Send(reply)
		return err
	}
}

//func ViewCmdAddSource(storage SourceStorage) botkit.ViewFunc {
//	type addSourceArgs struct {
//		Name     string `json:"name"`
//		URL      string `json:"url"`
//		Priority int    `json:"priority"`
//	}
//
//	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
//		args, err := botkit.ParseJSON[addSourceArgs](update.Message.CommandArguments())
//		if err != nil {
//			return err
//			// TODO: send error message
//		}
//
//		source := model.Source{
//			Name:     args.Name,
//			FeedURL:  args.URL,
//			Priority: args.Priority,
//		}
//
//		sourceID, err := storage.Add(ctx, source)
//		if err != nil {
//			return err
//		}
//
//		var (
//			msgText = fmt.Sprintf(
//				"Источник добавлен с ID: `%d`\\. Используйте его для управление источником\\.",
//				sourceID,
//			)
//			reply = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
//		)
//
//		reply.ParseMode = "MarkdownV2"
//
//		if _, err := bot.Send(reply); err != nil {
//			return err
//		}
//
//		return nil
//	}
//}
