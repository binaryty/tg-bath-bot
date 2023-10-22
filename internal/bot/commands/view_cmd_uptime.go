package commands

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	"github.com/yellowpuki/tg-bath-bot/internal/lib/countdown"
)

func ViewCmdUptime(t time.Time) bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		reply := "Uptime:\n" + countdown.Countdown{}.Count(t).String()
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, reply)); err != nil {
			return err
		}
		return nil
	}
}
