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

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yellowpuki/tg-bath-bot/internal/bot"
	"github.com/yellowpuki/tg-bath-bot/internal/bot/commands"
	"github.com/yellowpuki/tg-bath-bot/internal/storage/mongo"
)

const (
	DBUrl          = "mongodb://localhost:27017"
	ConnectTimeout = 10
	BotHost        = "api.telegram.org"
	ChatId         = 6421080707
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

	bathBot := bot.New(botApi, storage)
	bathBot.RegisterCmd("start", commands.ViewCmdStart())
	bathBot.RegisterCmd("help", commands.ViewCmdHelp())
	bathBot.RegisterCmd("uptime", commands.ViewCmdUptime(StartTime))
	bathBot.RegisterCmd("reg", commands.ViewCmdReg(ctx, storage))
	bathBot.RegisterCmd("last", commands.ViewCmdLast(ctx, storage))

	if err := bathBot.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] exiting: %v", err)
			os.Exit(1)
		}

		log.Println("bot stopped")
	}

}

// mustToken gets the token from the command line argument.
func mustToken() string {
	token := flag.String("t", "", "token for access telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token must be specified")
	}

	return *token
}
