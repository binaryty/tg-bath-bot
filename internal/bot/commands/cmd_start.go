package commands

import (
	"context"
	"log"

	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// CmdStart displays information on the time from the start of the bot.
func CmdStart() bot.CmdFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		log.Println(update.Message.Chat.ID)

		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, MsgStart)); err != nil {
			return err
		}
		return nil
	}
}
