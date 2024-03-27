package main

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/news-bot/internal/bot"
	"github.com/k6mil6/news-bot/internal/bot/middleware"
	"github.com/k6mil6/news-bot/internal/botkit"
	"github.com/k6mil6/news-bot/internal/config"
	"github.com/k6mil6/news-bot/internal/fetcher"
	"github.com/k6mil6/news-bot/internal/notifier"
	"github.com/k6mil6/news-bot/internal/state"
	"github.com/k6mil6/news-bot/internal/state/handlers"
	"github.com/k6mil6/news-bot/internal/storage"
	"github.com/k6mil6/news-bot/internal/summary"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	if err != nil {
		log.Printf("[ERROR] failed to create bot api: %s", err)
		return
	}

	db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	if err != nil {
		log.Printf("[ERROR] failed to connect to db: %v", err)
		return
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("[ERROR] failed to close db: %v", err)
		}
	}(db)

	var (
		articleStorage = storage.NewArticleStorage(db)
		sourceStorage  = storage.NewSourceStorage(db)
		f              = fetcher.New(
			articleStorage,
			sourceStorage,
			config.Get().FetchInterval,
			config.Get().FilterKeywords,
		)
		//summariser = summary.NewOpenAISummariser(
		//	config.Get().OpenAIKey,
		//	config.Get().OpenAIModel,
		//	config.Get().OpenAIPrompt,
		//)

		summariser = summary.NewLocalSummariser(
			config.Get().HTTPServerURL,
		)
		n = notifier.New(
			articleStorage,
			summariser,
			botAPI,
			config.Get().NotificationInterval,
			2*config.Get().FetchInterval,
			config.Get().TelegramChannelID,
		)
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	stateMachine := state.NewMachine()

	newsBot := botkit.New(botAPI, stateMachine)
	newsBot.RegisterCmdView("start", bot.ViewCmdStart(stateMachine))
	newsBot.RegisterCmdView(
		"add_source",
		middleware.AdminOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdAddSource(stateMachine),
		),
	)
	newsBot.RegisterCmdView(
		"list_sources",
		middleware.AdminOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdListSources(sourceStorage),
		),
	)
	newsBot.RegisterCmdView(
		"delete_source",
		middleware.AdminOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdDeleteSource(sourceStorage),
		),
	)
	newsBot.RegisterCmdView(
		"get_source",
		middleware.AdminOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdGetSource(sourceStorage),
		),
	)
	newsBot.RegisterCmdView(
		"set_priority",
		middleware.AdminOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdSetPriority(stateMachine),
		),
	)
	newsBot.RegisterCmdView(
		"stop_notifying_for",
		middleware.AdminOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdStopNotifyingFor(n),
		),
	)
	newsBot.RegisterCmdView(
		"stop_notifying",
		middleware.AdminOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdStopNotifying(n),
		),
	)
	newsBot.RegisterCmdView(
		"start_notifying",
		middleware.AdminOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdStartNotifying(n),
		),
	)

	newsBot.RegisterStateView(state.WaitingForSourceName, handlers.ViewWaitingForSourceName(sourceStorage, stateMachine))
	newsBot.RegisterStateView(state.WaitingForSourceURL, handlers.ViewWaitingForSourceURL(sourceStorage, stateMachine))
	newsBot.RegisterStateView(state.WaitingForSourcePriority, handlers.ViewWaitingForSourcePriority(sourceStorage, stateMachine))
	newsBot.RegisterStateView(state.WaitingForSourceIDAndPriority, handlers.ViewWaitingForSourceIDAndPriority(sourceStorage))

	if err := bot.SetCommands(botAPI, config.Get().TelegramChannelID); err != nil {
		log.Printf("[ERROR] failed to set commands: %s", err)
		return
	}

	go func(ctx context.Context) {
		if err = f.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to start f: %s", err)
				return
			}

			log.Println("[INFO] f stopped")
		}
	}(ctx)

	go func(ctx context.Context) {
		if err = n.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to start n: %s", err)
				return
			}

			log.Println("[INFO] n stopped")
		}

	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		log.Printf("[ERROR] failed to run botkit: %v", err)
	}
}
