package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	"github.com/yellowpuki/tg-bath-bot/internal/bot/commands"
	"github.com/yellowpuki/tg-bath-bot/internal/bot/mw"
	"github.com/yellowpuki/tg-bath-bot/internal/storage/db"
	"github.com/yellowpuki/tg-bath-bot/internal/storage/mongo"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const (
	DBUrl          = "mongodb://localhost:27017"
	ConnectTimeout = 10
	BotHost        = "api.telegram.org"
	ChatId         = -1002060428320
)

var StartTime = time.Now()

func main() {
	token := mustToken()

	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("[ERROR] can't start bot: %v", err)
		os.Exit(1)
	}

	log.Printf("[INFO] Authorized on account %s", botApi.Self.UserName)

	log.Println("[INFO] Start processing messages")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	storage := mongo.New(ctx, DBUrl, ConnectTimeout)

	db, err := db.New()
	if err != nil {
		log.Printf("[ERROR] can't access article storage: %v", err)
		os.Exit(1)
	}

	bathBot := bot.New(botApi, db)
	bathBot.RegisterCmd("start", commands.CmdStart())
	bathBot.RegisterCmd("help", commands.CmdHelp())
	bathBot.RegisterCmd("uptime", commands.CmdUptime(StartTime))
	bathBot.RegisterCmd("reg", mw.AdmOnly(ChatId, commands.CmdReg(ctx, storage)))
	bathBot.RegisterCmd("last", commands.CmdLast(ctx, storage))
	bathBot.RegisterCmd("art", mw.AdmOnly(ChatId, commands.CmdArticles(db)))
	bathBot.RegisterCmd("rnd", commands.CmdRndArticle(db))

	if err := bathBot.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] exiting: %v", err)
			os.Exit(1)
		}

		log.Println("[INFO] bot stopped")
	}

}

// mustToken gets the token from the command line argument.
func mustToken() string {
	token := flag.String("t", "", "token for access telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("[FATAL ERROR] token must be specified")
	}

	return *token
}
