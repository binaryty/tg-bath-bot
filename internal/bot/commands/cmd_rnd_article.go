package commands

import (
	"context"
	"fmt"
	"log"

	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	"github.com/yellowpuki/tg-bath-bot/internal/storage/db"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// CmdRndArticle outputs a random article from saved in the database.
func CmdRndArticle(db *db.DB) bot.CmdFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		art, err := db.GetRndArticle()

		if err != nil {
			log.Printf("[ERROR] can't get last event: %v", err)
			return err
		}

		reply := fmt.Sprintf(
			"[%s](%s)\n"+
				"Дата публикации на habr: _%s_\n",
			art.Title,
			art.URL,
			art.PublishedAt)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ParseMode = "markdown"
		if _, err := bot.Send(msg); err != nil {
			return err
		}
		return nil
	}
}
