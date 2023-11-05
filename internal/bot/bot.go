package bot

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/binaryty/tg-bath-bot/internal/storage/db"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const bathSize = 50 // default size of bath inline query offset

// Bot structure.
type Bot struct {
	api     *tgbotapi.BotAPI
	cmdMenu map[string]CmdFunc
	db      *db.DB
}

type CmdFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error

// New create a new bot instance.
func New(api *tgbotapi.BotAPI, db *db.DB) *Bot {
	return &Bot{
		api: api,
		db:  db,
	}
}

// Run start bot.
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

// RegisterCmd register a command in the bot menu.
func (b *Bot) RegisterCmd(cmd string, cmdFunc CmdFunc) {
	if b.cmdMenu == nil {
		b.cmdMenu = make(map[string]CmdFunc)
	}

	b.cmdMenu[cmd] = cmdFunc
}

// HandleUpdate handle an updates from telegram.
func (b *Bot) HandleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()

	if update.Message == nil && update.InlineQuery != nil {
		if err := b.ProcessInlineQuery(update); err != nil {
			log.Printf("[ERROR] can't process inline query: %v", err)
			return
		}
		log.Printf("[INFO] got a new inline query: [from]: %s [subject]: %s", update.InlineQuery.From.UserName, update.InlineQuery.Query)
	} else {
		if update.Message == nil || !update.Message.IsCommand() {
			return
		}

		log.Printf("[INFO] got a new update: [from]: %s [subject]: %s", update.Message.From.UserName, update.Message.Text)

		cmd := update.Message.Command()

		cmdFunc, ok := b.cmdMenu[cmd]
		if !ok {
			return
		}

		if err := cmdFunc(ctx, b.api, update); err != nil {
			log.Printf("[ERROR] failed to handle update: %v", err)

			if _, err := b.api.Send(
				tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n%v", "[Внутренняя ошибка]", err)),
			); err != nil {
				log.Printf("[ERROR] failed to send message: %v", err)
			}
		}

	}
}

// ProcessInlineQuery processing inline query messages from telegram.
func (b *Bot) ProcessInlineQuery(update tgbotapi.Update) error {
	inlineQuery := update.InlineQuery
	queryOffset, _ := strconv.Atoi(inlineQuery.Offset)

	if queryOffset == 0 {
		queryOffset = 1
	}

	results := make([]interface{}, 0)

	articles, err := b.db.GetArticlesByTitle(strings.ToLower(inlineQuery.Query))
	if err != nil {
		return err
	}

	for _, article := range offsetResult(queryOffset, articles) {
		msg := fmt.Sprintf("%s %s", article.Title, article.URL)
		results = append(results, tgbotapi.InlineQueryResultArticle{
			Type:  "article",
			ID:    article.Id(),
			Title: article.Title,
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text: msg,
			},
			ThumbURL: article.ThumbURL,
		})
	}

	if len(results) < 50 {
		_, err := b.api.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID: inlineQuery.ID,
			Results:       results,
			IsPersonal:    true,
			CacheTime:     0,
		})

		if err != nil {
			return err
		}
	} else {
		_, err := b.api.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID: inlineQuery.ID,
			Results:       results,
			IsPersonal:    true,
			CacheTime:     0,
			NextOffset:    strconv.Itoa(queryOffset + bathSize),
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// offsetResult ...
func offsetResult(startNum int, articles []db.Article) []db.Article {
	overallItems := len(articles)

	switch {
	case startNum >= overallItems:
		return []db.Article{}
	case startNum+bathSize >= overallItems:
		return articles[startNum:overallItems]
	default:
		return articles[startNum : startNum+bathSize]
	}
}
