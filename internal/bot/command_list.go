package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var adminCommands = []tgbotapi.BotCommand{
	{
		Command:     "list_sources",
		Description: "Команда для получения списка источников",
	},
	{
		Command:     "delete_source",
		Description: "Команда для удаления источника, по его идентификатору (/delete_source <идентификатор>)",
	},
	{
		Command:     "get_source",
		Description: "Команда для получения источника по его идентификатору (/get_source <идентификатор>)",
	},
	{
		Command:     "set_priority",
		Description: `Команда для установки приоритета`,
	},
	//{
	//	Command:     "stop_notifying_for",
	//	Description: "Команда для остановки постинга статей (/stop_notifying_for <длительность>)",
	//},
	{
		Command:     "add_source",
		Description: `Команда для добавления источника (/add_source {"name": <имя>, ""url": <адрес>, "priority": <приоритет>})`,
	},
	{
		Command:     "stop_notifying",
		Description: `Команда для остановки постинга статей`,
	},
	{
		Command:     "start_notifying",
		Description: `Команда для запуска постинга статей`,
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
