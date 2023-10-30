package commands

import (
	"context"
	"log"

	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func ViewCmdStart() bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		log.Println(update.Message.Chat.ID)

		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, MsgStart)); err != nil {
			return err
		}
		return nil
	}
}
