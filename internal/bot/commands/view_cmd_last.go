package commands

import (
	"context"
	"log"

	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	"github.com/yellowpuki/tg-bath-bot/internal/lib/countdown"
	"github.com/yellowpuki/tg-bath-bot/internal/storage"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func ViewCmdLast(ctx context.Context, s storage.Storage) bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		last, err := s.LastVisit(ctx, update.Message.Chat.UserName)

		if err != nil {
			log.Printf("[ERROR] can't get last event: %v", err)
			return err
		}

		reply := "Last event is:\n" + countdown.Countdown{}.Count(last).String()

		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, reply)); err != nil {
			return err
		}
		return nil
	}
}
