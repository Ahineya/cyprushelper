package main

import (
	"os"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"github.com/Ahineya/cyprushelper/commandrouter"
	"github.com/Ahineya/cyprushelper/services/pollutionservice"
	"github.com/Ahineya/cyprushelper/helpers/logger"
)

func main() {
	if os.Getenv("ENV") == "PROD" {
		runProd()

	} else {
		runDev()
	}
}

func runProd() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		logger.Error("MAIN", "You need to set PORT environment variable")
		os.Exit(1)
	}

	token := os.Getenv("BOT_TOKEN")
	if len(token) == 0 {
		logger.Error("MAIN", "You need to set BOT_TOKEN environment variable")
		os.Exit(1)
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	if len(webhookURL) == 0 {
		logger.Error("MAIN", "You need to set WEBHOOK_URL environment variable")
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Error("MAIN", err.Error())
		os.Exit(1)
	}

	logger.Info("BOT", "Authorized on account " + bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL + bot.Token))
	if err != nil {
		logger.Error("MAIN", err.Error())
		os.Exit(1)
	}

	bot.Debug = true

	setupServices(bot)

	updates := bot.ListenForWebhook("/" + bot.Token)

	go http.ListenAndServe("0.0.0.0:" + port, nil)

	for update := range updates {
		logger.Log("BOT", update.Message.From.UserName + " " + update.Message.Text)

		if (len(update.Message.Text) > 1 && string(update.Message.Text[0]) == "/") {
			commandrouter.Route(bot, update)
		}
	}
}

func runDev() {
	token := os.Getenv("BOT_TOKEN")
	if len(token) == 0 {
		logger.Error("MAIN", "You need to set BOT_TOKEN environment variable")
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Error("MAIN", err.Error())
		os.Exit(1)
	}

	bot.Debug = true

	logger.Info("BOT", "Authorized on account " + bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	setupServices(bot)

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		logger.Log("BOT", update.Message.From.UserName + " " + update.Message.Text)

		if (len(update.Message.Text) > 1 && string(update.Message.Text[0]) == "/") {
			commandrouter.Route(bot, update)
		}
	}
}

func setupServices(bot *tgbotapi.BotAPI) {
	pollutionservice.Start(bot)
}

