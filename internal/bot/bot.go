package bot

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/yellowpuki/tg-bath-bot/internal/lib/er"
	"github.com/yellowpuki/tg-bath-bot/internal/storage/db"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const queryLimit = 200

type Bot struct {
	api     *tgbotapi.BotAPI
	cmdView map[string]ViewFunc
	db      *db.DB
}

type ViewFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

func New(api *tgbotapi.BotAPI, db *db.DB) *Bot {
	return &Bot{
		api: api,
		db:  db,
	}
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.api.GetUpdatesChan(u)
	if err != nil {
		return err
	}

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
			log.Printf("[ERROR] panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()

	if update.Message == nil && update.InlineQuery != nil {
		b.handleInlineQuery(update)
		//if err := b.ProcessInlineQuery(update); err != nil {
		//	log.Printf("[ERROR] can't process inline query: %v", err)
		//	return
		//}
		log.Printf("[INFO] got a new inline query: [from]: %s [subject]: %s", update.InlineQuery.From.UserName, update.InlineQuery.Query)
	} else {
		if update.Message == nil || !update.Message.IsCommand() {
			return
		}

		log.Printf("[INFO] got a new update: [from]: %s [subject]: %s", update.Message.From.UserName, update.Message.Text)

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
}

func (b *Bot) ProcessInlineQuery(update tgbotapi.Update) error {
	query := update.InlineQuery.Query
	queryOffset, _ := strconv.Atoi(update.InlineQuery.Offset)
	if queryOffset == 0 {
		queryOffset = 1
	}

	articles, err := b.db.GetArticlesByTitle(strings.ToLower(query), queryLimit)
	if err != nil {
		return er.Wrap("get articles", err)
	}

	var results []interface{}

	for _, article := range articles {
		msg := fmt.Sprintf("%s %s", article.Title, article.URL)
		res := tgbotapi.NewInlineQueryResultArticle(article.Id(), article.Title, msg)
		results = append(results, res)
	}

	var inlineConf tgbotapi.InlineConfig

	inlineConf = tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       results,
	}
	_, err = b.api.AnswerInlineQuery(inlineConf)
	if err != nil {
		return er.Wrap("answer inline query", err)
	}

	return nil
}

func (b *Bot) handleInlineQuery(update tgbotapi.Update) {
	inlineQuery := update.InlineQuery
	offset, _ := strconv.Atoi(inlineQuery.Offset)

	results := make([]interface{}, 0)

	articles, _ := b.db.GetArticlesByTitle(strings.ToLower(inlineQuery.Query), queryLimit)

	for _, article := range articles {
		inlineResult := tgbotapi.NewInlineQueryResultArticle(article.Id(), article.Title, article.Title)
		results = append(results, inlineResult)
	}

	if len(results) < 50 {
		inlineConfig := tgbotapi.InlineConfig{
			InlineQueryID: inlineQuery.ID,
			Results:       results,
			IsPersonal:    true,
			CacheTime:     0,
		}
		_, err := b.api.AnswerInlineQuery(inlineConfig)
		if err != nil {
			log.Println(err)
		}
	} else {
		inlineConfig := tgbotapi.InlineConfig{
			InlineQueryID: inlineQuery.ID,
			Results:       results,
			IsPersonal:    true,
			CacheTime:     0,
			NextOffset:    strconv.Itoa(offset + 50),
		}
		_, err := b.api.AnswerInlineQuery(inlineConfig)
		if err != nil {
			log.Println(err)
		}
	}
}
