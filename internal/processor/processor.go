package processor

import (
	"context"
	"errors"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yellowpuki/tg-bath-bot/internal/countdown"
	"github.com/yellowpuki/tg-bath-bot/internal/er"
	"github.com/yellowpuki/tg-bath-bot/internal/storage"
	"golang.org/x/exp/slog"
)

var ErrUserExists = errors.New("user exists")

type Processor struct {
	client  *tgbotapi.BotAPI
	storage storage.Storage
}

// New ...
func New(client *tgbotapi.BotAPI, storage storage.Storage) *Processor {

	return &Processor{
		client:  client,
		storage: storage,
	}
}

// GetLastDate ...
func (p *Processor) LastDateCmd(ctx context.Context, upd tgbotapi.Update) tgbotapi.MessageConfig {
	var reply string

	last, err := p.storage.LastVisit(ctx, upd.Message.Chat.UserName)

	if err != nil {
		slog.Error("can't get", slog.String("Error", err.Error()))
		return tgbotapi.NewMessage(upd.Message.Chat.ID, err.Error())
	}

	reply = "Last bath is:\n" + countdown.Countdown{}.Count(last).String()

	return tgbotapi.NewMessage(upd.Message.Chat.ID, reply)

}

// Register ...
func (p *Processor) RegisterCmd(ctx context.Context, upd tgbotapi.Update, r *storage.Record) tgbotapi.MessageConfig {
	user := upd.Message.From.UserName
	err := reg(context.Background(), p.storage, r)
	if err != nil {
		slog.Error("can't register user", slog.String("Error", err.Error()))

		return tgbotapi.NewMessage(upd.Message.From.ID, err.Error())
	}

	reply := fmt.Sprintf("[%s] succesfully register", user)
	return tgbotapi.NewMessage(upd.Message.Chat.ID, reply)

}

// StartCmd
func (p *Processor) StartCmd(upd tgbotapi.Update) tgbotapi.MessageConfig {

	return tgbotapi.NewMessage(upd.Message.Chat.ID, msgStart)

}

// HelpCmd ...
func (p *Processor) HelpCmd(upd tgbotapi.Update) tgbotapi.MessageConfig {

	return tgbotapi.NewMessage(upd.Message.Chat.ID, msgHelp)
}

// Uptime ...
func (p *Processor) UptimeCmd(upd tgbotapi.Update, t time.Time) tgbotapi.MessageConfig {

	return tgbotapi.NewMessage(upd.Message.Chat.ID, uptime(t))
}

// MenuCmd ...
func (p *Processor) MenuCmd(upd tgbotapi.Update) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Select command")

	keyboard := tgbotapi.InlineKeyboardMarkup{}

	for _, command := range commands {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(command, command)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg.ReplyMarkup = keyboard

	return &msg
}

// uptime ...
func uptime(t time.Time) string {

	return "Uptime:\n" + countdown.Countdown{}.Count(t).String()
}

// reg ...
func reg(ctx context.Context, s storage.Storage, r *storage.Record) error {
	op := "main:Register"

	h, err := r.EventHash()
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
		slog.Info("register: record exists", slog.String("user", r.UserName))
		return er.Wrap(op, ErrUserExists)
	}

	rec := &storage.Record{
		EventToken: h,
		UserName:   r.UserName,
		CreatedAt:  time.Now(),
	}
	err = s.Save(context.Background(), rec)

	if err != nil {
		return er.Wrap("can't register", err)
	}

	return nil
}
