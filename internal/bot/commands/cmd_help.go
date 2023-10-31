package commands

import (
	"context"

	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// CmdHelp displays a help information.
func CmdHelp() bot.CmdFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, MsgHelp)); err != nil {
			return err
		}
		return nil
	}
}
