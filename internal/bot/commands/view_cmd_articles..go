package commands

import (
	"context"
	"errors"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	"github.com/yellowpuki/tg-bath-bot/internal/lib/fetcher"
	"github.com/yellowpuki/tg-bath-bot/internal/storage/db"
)

func ViewCmdArticles(s *db.DB) bot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		f := fetcher.New(500 * time.Second)

		go func(ctx context.Context) {
			if err := f.Start(ctx); err != nil {
				if !errors.Is(err, context.Canceled) {
					log.Printf("[ERROR] failed to run fetcher: %v", err)
					return
				}
				log.Printf("[INFO] fetcher stopped")

				for _, a := range f.Articles {
					if err := s.SaveArticle(a); err != nil {
						log.Printf("[ERROR] can't save article: %v", err)
						return
					}
				}
			}
		}(ctx)

		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, msgSaved)); err != nil {
			return err
		}
		return nil
	}
}
