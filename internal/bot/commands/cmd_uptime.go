package commands

import (
	"context"
	"time"

	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	"github.com/yellowpuki/tg-bath-bot/internal/lib/countdown"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// CmdUpteime displays information on the time from the start of the bot.
func CmdUptime(t time.Time) bot.CmdFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		reply := countdown.Countdown{}.Count(t).String()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		msg.ParseMode = "markdown"
		if _, err := bot.Send(msg); err != nil {
			return err
		}
		return nil
	}
}
