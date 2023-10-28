package commands

import (
	"context"
	"fmt"
	"log"

	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	"github.com/yellowpuki/tg-bath-bot/internal/storage/db"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func ViewCmdRndArticle(db *db.DB) bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		art, err := db.GetRndArticle()

		if err != nil {
			log.Printf("[ERROR] can't get last event: %v", err)
			return err
		}

		reply := fmt.Sprintf("%s, %s", art.Title, art.URL)

		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, reply)); err != nil {
			return err
		}
		return nil
	}
}
