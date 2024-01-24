package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var adminCommands = []tgbotapi.BotCommand{
	{
		Command:     "listsources",
		Description: "Команда для получения списка источников",
	},
	{
		Command:     "deletesource",
		Description: "Команда для удаления источника, по его идентификатору (/deletesource <идентификатор>)",
	},
	{
		Command:     "getsource",
		Description: "Команда для получения источника по его идентификатору (/getsource <идентификатор>)",
	},
	{
		Command:     "setpriority",
		Description: `Команда для установки приоритета (/setpriority {"source_id": <идентификатор>, "priority": <приоритет>})`,
	},
	{
		Command:     "stopnotifyingfor",
		Description: "Команда для остановки постинга статей (/stopnotifyingfor <длительность>)",
	},
	{
		Command:     "addsource",
		Description: `Команда для добавления источника (/addsource {"name": <имя>, ""url": <адрес>, "priority": <приоритет>})`,
	},
}

func SetCommands(bot *tgbotapi.BotAPI, channelID int64) error {
	admins, err := bot.GetChatAdministrators(
		tgbotapi.ChatAdministratorsConfig{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: channelID,
			},
		},
	)

	if err != nil {
		return err
	}

	for _, admin := range admins {
		scope := tgbotapi.NewBotCommandScopeChat(admin.User.ID)
		setCommands := tgbotapi.SetMyCommandsConfig{
			Commands: adminCommands,
			Scope:    &scope,
		}

		if _, err := bot.Request(setCommands); err != nil {
			return err
		}
	}

	return nil
}
