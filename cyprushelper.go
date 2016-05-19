package main

import (
	"os"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"github.com/Ahineya/cyprushelper/pollution"
	"github.com/Ahineya/cyprushelper/storage"
	"github.com/Ahineya/cyprushelper/commandrouter"
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
		panic("You need to set PORT environment variable")
	}

	token := os.Getenv("BOT_TOKEN")
	if len(token) == 0 {
		panic("You need to set BOT_TOKEN environment variable")
	}

	webhookURL := os.Getenv("WEBHOOK_URL")
	if len(webhookURL) == 0 {
		panic("You need to set WEBHOOK_URL environment variable")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL + bot.Token))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	// Temporarily disable services for prod
	//setupServices(bot)

	updates := bot.ListenForWebhook("/" + bot.Token)

	go http.ListenAndServe("0.0.0.0:" + port, nil)

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if (len(update.Message.Text) > 1 && string(update.Message.Text[0]) == "/") {
			commandrouter.Route(bot, update)
		}
	}
}

func runDev() {
	token := os.Getenv("BOT_TOKEN")
	if len(token) == 0 {
		panic("You need to set BOT_TOKEN environment variable")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	setupServices(bot)

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if (len(update.Message.Text) > 1 && string(update.Message.Text[0]) == "/") {
			commandrouter.Route(bot, update)
		}
	}
}

func setupServices(bot *tgbotapi.BotAPI) {
	// Setting up a pollution service
	pollutionServiceChannel := make(chan string)
	pollution.CreatePollutionService(pollutionServiceChannel)
	go func() {
		for pollutionData := range pollutionServiceChannel {

			// Get all chat Ids
			chats, err := storage.GetChatIds()
			if err != nil {
				continue
			}
			// Send messages to all chats
			for _, chatId := range chats {
				msg := tgbotapi.NewMessage(chatId, pollutionData)
				bot.Send(msg)
			}

			// TODO: Set a ticker for batch sending updates to ids

		}
	}()
}

