package commands

import (
	"context"
	"time"

	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	"github.com/yellowpuki/tg-bath-bot/internal/lib/countdown"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func ViewCmdUptime(t time.Time) bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		reply := "Uptime:\n" + countdown.Countdown{}.Count(t).String()
		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, reply)); err != nil {
			return err
		}
		return nil
	}
}
