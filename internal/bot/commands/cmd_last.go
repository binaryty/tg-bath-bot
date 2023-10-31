package commands

import (
	"context"
	"log"

	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	"github.com/yellowpuki/tg-bath-bot/internal/lib/countdown"
	"github.com/yellowpuki/tg-bath-bot/internal/storage"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// CmdLast displays information on the time tracking of the last event.
func CmdLast(s storage.Storage) bot.CmdFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		last, err := s.LastVisit(ctx, update.Message.From.UserName)

		if err != nil {
			log.Printf("[ERROR] can't get last event: %v", err)
			return err
		}

		reply := "Last event is:\n" + countdown.Countdown{}.Count(last).String()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ParseMode = "markdown"
		if _, err := bot.Send(msg); err != nil {
			return err
		}
		return nil
	}
}
