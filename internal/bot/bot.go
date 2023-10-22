package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yellowpuki/tg-bath-bot/internal/storage"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	cmdView map[string]ViewFunc
	Storage storage.Storage
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

func New(api *tgbotapi.BotAPI, storage storage.Storage) *Bot {
	return &Bot{
		api:     api,
		Storage: storage,
	}
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(ctx, 5*time.Second)
			b.HandleUpdate(updateCtx, update)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) RegisterCmd(cmd string, view ViewFunc) {
	if b.cmdView == nil {
		b.cmdView = make(map[string]ViewFunc, 0)
	}

	b.cmdView[cmd] = view
}

func (b *Bot) HandleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered: %v", p)
		}
	}()

	if update.Message == nil || !update.Message.IsCommand() {
		return
	}

	cmd := update.Message.Command()

	view, ok := b.cmdView[cmd]
	if !ok {
		return
	}

	if err := view(ctx, b.api, update); err != nil {
		log.Printf("[ERROR] failed to handle update: %v", err)

		if _, err := b.api.Send(
			tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s:%v", "internal error", err)),
		); err != nil {
			log.Printf("[ERROR] failed to send message: %v", err)
		}
	}

}
