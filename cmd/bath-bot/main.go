package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yellowpuki/tg-bath-bot/internal/processor"
	stor "github.com/yellowpuki/tg-bath-bot/internal/storage"
	"github.com/yellowpuki/tg-bath-bot/internal/storage/mongo"
	"golang.org/x/exp/slog"
)

const (
	DBUrl          = "mongodb://localhost:27017"
	ConnectTimeout = 10
	BotHost        = "api.telegram.org"
)

var appStartTime = time.Now()

func main() {
	token := mustToken()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	ctx := context.Background()

	storage := mongo.New(ctx, DBUrl, ConnectTimeout)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Error("can't start bot: %s", err)
		os.Exit(1)
	}

	log.Info("Authorized on account", slog.String("Account", bot.Self.UserName))

	p := processor.New(bot, storage)

	log.Info("Start processing messages")

	updConfig := tgbotapi.NewUpdate(0)
	updConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updConfig)

	commandMap := make(map[int64]string)

	for update := range updates {
		if update.Message != nil {
			log.Info("New message:", slog.String("from", update.Message.From.UserName), slog.String("message", update.Message.Text))

			command := update.Message.Command()

			userName := update.Message.Chat.UserName

			switch command {
			case "start":
				bot.Send(p.StartCmd(update))
			case "help":
				bot.Send(p.HelpCmd(update))
			case "uptime":
				bot.Send(p.UptimeCmd(update, appStartTime))
			case "reg":
				bot.Send(p.RegisterCmd(ctx, update, &stor.Record{UserName: userName}))
			case "last":
				bot.Send(p.LastDateCmd(ctx, update))
			case "menu":
				bot.Send(p.MenuCmd(update))
			}
		} else {
			if update.CallbackQuery != nil {
				c := update.CallbackQuery.Data
				commandMap[update.CallbackQuery.From.ID] = c
				bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ok"))
			}

		}
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
