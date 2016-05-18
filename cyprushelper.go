package main

import (
	"os"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"github.com/Ahineya/cyprushelper/pharmacies"
	//"github.com/Ahineya/cyprushelper/bypass-asp"
	//"fmt"
	"github.com/Ahineya/cyprushelper/seatemp"
	"net/http"
	"github.com/Ahineya/cyprushelper/stats"
	"github.com/Ahineya/cyprushelper/pollution"
)

var Messages = map[string]string{
	"Start": "Hi! I am Cyprus Helper bot!",
	"Help": "You can use the following commands: /pharmacies",
}

const (
	// TODO: define as environment variable
	bot_name = "CyprusHelper_bot"
)

func main() {
	if os.Getenv("ENV") == "PROD" {
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

		updates := bot.ListenForWebhook("/" + bot.Token)

		go http.ListenAndServe("0.0.0.0:" + port, nil)

		for update := range updates {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if (len(update.Message.Text) > 1 && string(update.Message.Text[0]) == "/") {
				processUpdate(bot, update)
			}
		}

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
	if os.Getenv("ENV") == "PROD" {
		stats.Track(update.Message)
	}

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
	case "/seatemp":
		sendSeaTemp(bot, update, tokens)
	case "/seatemp" + bot_name:
		sendSeaTemp(bot, update, tokens)
	case "/pollution":
		sendPollution(bot, update, tokens)
	case "/pollution" + bot_name:
		sendPollution(bot, update, tokens)
	}
}

func sendPharmacies(bot *tgbotapi.BotAPI, update tgbotapi.Update, tokens []string) {
	allPharmacies, err := pharmacies.GetAllPharmacies()
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Some error happen: " + err.Error())
		bot.Send(msg)
		return;
	}

	// Send info for all cities
	if len(tokens) < 2 {
		for _, pharmaciesInCity := range allPharmacies {
			text := pharmacies.FormatPharmacies(pharmaciesInCity)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join([]string{text}, "\n"))
			msg.ParseMode = tgbotapi.ModeHTML
			bot.Send(msg)
		}
		return;
	}

	// If envoked as /pharmacies CITY_NAME
	city := pharmacies.NormalizeCity(strings.ToLower(tokens[1]))
	found := false
	for _, pharmaciesInCity := range allPharmacies {
		if city == pharmaciesInCity.City {
			found = true
			text := pharmacies.FormatPharmacies(pharmaciesInCity)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join([]string{text}, "\n"))
			msg.ParseMode = tgbotapi.ModeHTML
			bot.Send(msg)
			break;
		}
	}

	if !found {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "City " + `"` + city + `" not found`)
		bot.Send(msg)
	}

}

func sendSeaTemp(bot *tgbotapi.BotAPI, update tgbotapi.Update, tokens []string) {

	if len(tokens) < 2 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please, define city: /seatemp city_name")
		bot.Send(msg)
		return
	}

	city := strings.ToLower(tokens[1])
	seatemp, err := seatemp.GetSeaTemp(city)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Can't find a sea temperature for " + strings.Join(tokens[1:], " ") + ", error happen: " + err.Error())
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sea temperatire in " + strings.Join(tokens[1:], " ") +" is " + seatemp)
	bot.Send(msg)

}

func sendPollution(bot *tgbotapi.BotAPI, update tgbotapi.Update, tokens []string) {
	if len(tokens) < 2 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please, define city: /pollution city_name")
		bot.Send(msg)
		return
	}

	city := strings.ToLower(tokens[1])
	pollutionResult, err := pollution.GetPollution(city)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Can't find pollution info for " + strings.Join(tokens[1:], " ") + ", error happen: " + err.Error())
		bot.Send(msg)
		return;
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Pollution in " + strings.Join(tokens[1:], " ") + " is:\n\n" + pollution.FormatPollution(pollutionResult))
	msg.ParseMode = tgbotapi.ModeHTML
	bot.Send(msg)
}