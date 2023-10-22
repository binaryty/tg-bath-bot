package middleware

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yellowpuki/tg-bath-bot/internal/bot"
)

func RoleCheck(chanId int64, next bot.ViewFunc) bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		adms, err := bot.GetChatAdministrators(
			tgbotapi.ChatAdministratorsConfig{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: chanId,
				},
			})
		if err != nil {
			return err
		}

		for _, adm := range adms {
			if adm.User.ID == update.Message.From.ID {
				return next(ctx, bot, update)
			}
		}

		if _, err := bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"Access denide")); err != nil {
			return err
		}

		return nil
	}
}
