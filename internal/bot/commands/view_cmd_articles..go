package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	"github.com/yellowpuki/tg-bath-bot/internal/lib/fetcher"
	"github.com/yellowpuki/tg-bath-bot/internal/storage/db"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
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
						errMsg := fmt.Sprintf("[ERROR] can't save article: %v", err)
						log.Println(errMsg)
						if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, errMsg)); err != nil {
							return
						}
						return
					}
				}
			}
			if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msgSaved)); err != nil {
				return
			}
		}(ctx)

		return nil
	}
}
