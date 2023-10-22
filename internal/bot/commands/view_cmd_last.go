package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	"github.com/yellowpuki/tg-bath-bot/internal/lib/countdown"
	"github.com/yellowpuki/tg-bath-bot/internal/storage"
	"golang.org/x/exp/slog"
)

func ViewCmdLast(ctx context.Context, s storage.Storage) bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		last, err := s.LastVisit(ctx, update.Message.Chat.UserName)

		if err != nil {
			slog.Error("can't get", slog.String("Error", err.Error()))
			return err
		}

		reply := "Last bath is:\n" + countdown.Countdown{}.Count(last).String()

		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, reply)); err != nil {
			return err
		}
		return nil
	}
}
