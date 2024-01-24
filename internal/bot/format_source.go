package bot

import (
	"fmt"
	"github.com/k6mil6/news-bot/internal/botkit/markup"
	"github.com/k6mil6/news-bot/internal/model"
)

func formatSource(source model.Source) string {
	return fmt.Sprintf(
		"🌐 *%s*\nID: `%d`\nURL фида: %s\nПриоритет: %d",
		markup.EscapeForMarkdown(source.Name),
		source.ID,
		markup.EscapeForMarkdown(source.FeedURL),
		source.Priority,
	)
}
