package commands

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	"github.com/yellowpuki/tg-bath-bot/internal/lib/er"
	"github.com/yellowpuki/tg-bath-bot/internal/storage"
	"golang.org/x/exp/slog"
)

func ViewCmdReg(ctx context.Context, s storage.Storage) bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

		if err := reg(ctx, s, update); err != nil {
			return err
		}

		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, MsgReg)); err != nil {
			return err
		}
		return nil
	}
}

// reg ...
func reg(ctx context.Context, s storage.Storage, update tgbotapi.Update) error {

	user := update.Message.From.UserName

	h, err := storage.EventHash(user)
	if err != nil {
		return er.Wrap("can't register", err)
	}

	isExist, err := s.IsExist(ctx, h)
	if err != nil {
		if err != storage.ErrNoRecords {
			return er.Wrap("can't register", err)
		}
	}

	if isExist {
		slog.Info("register: record exists", slog.String("user", user))
		return er.ErrUserExists
	}

	rec := &storage.Record{
		EventToken: h,
		UserName:   user,
		CreatedAt:  time.Now(),
	}
	err = s.Save(ctx, rec)

	if err != nil {
		return er.Wrap("can't register", err)
	}

	return nil
}
