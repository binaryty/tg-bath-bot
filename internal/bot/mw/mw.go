package mw

import (
	"context"

	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// AdmOnly middleware to restrict access only for group admins.
func AdmOnly(chatId int64, next bot.CmdFunc) bot.CmdFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		admins, err := bot.GetChatAdministrators(
			tgbotapi.ChatConfig{
				ChatID: chatId,
			},
		)

		if err != nil {
			return err
		}

		for _, adm := range admins {
			if adm.User.ID == sentFrom(update).ID {
				return next(ctx, bot, update)
			}
		}

		if _, err := bot.Send(tgbotapi.NewMessage(
			fromChat(update).ID,
			"У вас нет прав на выполнение этой команды.",
		)); err != nil {
			return err
		}

		return nil
	}
}

func sentFrom(u tgbotapi.Update) *tgbotapi.User {
	switch {
	case u.Message != nil:
		return u.Message.From
	case u.EditedMessage != nil:
		return u.EditedMessage.From
	case u.InlineQuery != nil:
		return u.InlineQuery.From
	case u.ChosenInlineResult != nil:
		return u.ChosenInlineResult.From
	case u.CallbackQuery != nil:
		return u.CallbackQuery.From
	case u.ShippingQuery != nil:
		return u.ShippingQuery.From
	case u.PreCheckoutQuery != nil:
		return u.PreCheckoutQuery.From
	default:
		return nil

	}
}

func fromChat(u tgbotapi.Update) *tgbotapi.Chat {
	switch {
	case u.Message != nil:
		return u.Message.Chat
	case u.EditedMessage != nil:
		return u.EditedMessage.Chat
	case u.ChannelPost != nil:
		return u.ChannelPost.Chat
	case u.EditedChannelPost != nil:
		return u.EditedChannelPost.Chat
	case u.CallbackQuery != nil:
		return u.CallbackQuery.Message.Chat
	default:
		return nil
	}
}
