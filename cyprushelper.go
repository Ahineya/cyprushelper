package main

import (
	"os"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"github.com/Ahineya/cyprushelper/pharmacies"
	//"github.com/Ahineya/cyprushelper/bypass-asp"
	//"fmt"
)

var Messages = map[string]string{
	"Start": "Hi! I am Cyprus Helper bot!",
	"Help": "You can use the following commands: /pharmacies",
}

const (
	bot_name = "CyprusHelper_bot"
)

func main() {
	if os.Getenv("ENV") == "PROD" {
		panic("PROD not implemented")
		/*
			Here we will add a production version handler with webhook,
			like in https://github.com/go-telegram-bot-api/telegram-bot-api
			
			I think it will be good to have WEBHOOK_URL environment variable for that,
			to not expose it on the Github.
		 */
	} else {
		token := os.Getenv("BOT_TOKEN")
		if len(token) == 0 {
			panic("You need to set BOT_TOKEN environment variable")
		}
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Panic(err)
		}

		bot.Debug = true

		log.Printf("Authorized on account %s", bot.Self.UserName)

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err := bot.GetUpdatesChan(u)

		for update := range updates {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if (len(update.Message.Text) > 1 && string(update.Message.Text[0]) == "/") {
				processUpdate(bot, update)
			}
		}
	}
}

// Refactor all this switch stuff to different modules
func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	tokens := strings.Fields(update.Message.Text)
	command := tokens[0]

	switch command {
	case "/start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, Messages["Start"])
		bot.Send(msg)
	case "/start@" + bot_name:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, Messages["Start"])
		bot.Send(msg)
	case "/help":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, Messages["Help"])
		bot.Send(msg)
	case "/help@" + bot_name:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, Messages["Help"])
		bot.Send(msg)
	case "/pharmacies":
		sendPharmacies(bot, update, tokens)
	case "/pharmacies" + bot_name:
		sendPharmacies(bot, update, tokens)
	}
}

func sendPharmacies(bot *tgbotapi.BotAPI, update tgbotapi.Update, tokens []string) {
	allPharmacies, err := pharmacies.GetAllPharmacies()
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Some error happen: " + err.Error())
		bot.Send(msg)
	}

	// If envoked as /pharmacies CITY_NAME
	if len(tokens) > 1 {
		city := tokens[1]
		found := false
		for _, pharmaciesInCity := range allPharmacies {
			if city == pharmaciesInCity.City {
				found = true
				text := pharmacies.FormatPharmacies(pharmaciesInCity)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join([]string{text}, "\n"))
				msg.ParseMode = tgbotapi.ModeHTML
				bot.Send(msg)
			}
		}

		if !found {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "City " + `"` + city + `" not found`)
			bot.Send(msg)
		}
	} else {
		for _, pharmaciesInCity := range allPharmacies {
			text := pharmacies.FormatPharmacies(pharmaciesInCity)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join([]string{text}, "\n"))
			msg.ParseMode = tgbotapi.ModeHTML
			bot.Send(msg)
		}
	}
}


