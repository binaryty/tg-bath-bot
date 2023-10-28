package commands

import (
	"context"

	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func ViewCmdHelp() bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, MsgHelp)); err != nil {
			return err
		}
		return nil
	}
}
