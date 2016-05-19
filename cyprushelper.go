package main

import (
	"os"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"github.com/Ahineya/cyprushelper/pharmacies"
	"github.com/Ahineya/cyprushelper/seatemp"
	"net/http"
	"github.com/Ahineya/cyprushelper/stats"
	"github.com/Ahineya/cyprushelper/pollution"
	"github.com/Ahineya/cyprushelper/storage"
)

var Messages = map[string]string{
	"Start": "Hi! I am Cyprus Helper bot!",
	"Help": "You can use the following commands: /pharmacies",
}

const (
	// TODO: define as environment variable
	bot_name = "CyprusHelper_bot"

	// TODO: find proper place for such messages
	seatemp_commands_message = `
	Please, define city: /seatemp city_name
	Or use commands for available cities:

	Limassol: /seatemp_limassol
	Paphos: /seatemp_pafos
	Larnaka: /seatemp_larnaka
	Nicosia (Kato Pyrgos): /seatemp_nicosia
	Famagusta: /seatemp_famagusta
	Protaras: /seatemp_protaras
	Ayia-napa: /seatemp_ayianapa
	`

	pharmacies_commands_message = `
	Please, define city: /pharmacies city_name
	Or use commands for available cities:

	Limassol: /pharmacies_limassol
	Paphos: /pharmacies_pafos
	Larnaka: /pharmacies_larnaka
	Nicosia (Kato Pyrgos): /pharmacies_nicosia
	Famagusta: /pharmacies_famagusta
	`

	pollution_commands_message = `
	Please, define city: /pollution city_name
	Or use commands for available cities:

	Limassol: /pollution_limassol
	Paphos: /pollution_pafos
	Larnaka: /pollution_larnaka
	Nicosia (Kato Pyrgos): /pollution_nicosia
	`
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

		// Temporarily disable services for prod
		//setupServices(bot)

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

		setupServices(bot)

		updates, err := bot.GetUpdatesChan(u)

		for update := range updates {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if (len(update.Message.Text) > 1 && string(update.Message.Text[0]) == "/") {
				processUpdate(bot, update)
			}
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

// Refactor all this switch stuff to different modules
func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	//if os.Getenv("ENV") == "PROD" {
		stats.Track(update.Message)
		go storage.UpdateChats(update.Message)
	//}

	tokens := strings.Fields(update.Message.Text)
	command := tokens[0]

	// Checking for composite cxommands like /pharmacies_limassol instead of /pharmacies limassol
	compositeCommand := strings.Split(tokens[0], "_")

	if len(compositeCommand) == 2 {
		command = compositeCommand[0]
		if len(tokens) == 1 {
			tokens = append(tokens, compositeCommand[1])
		} else {
			tokens[1] = compositeCommand[1]
		}
	}

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

	if len(tokens) < 2 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, pharmacies_commands_message)
		bot.Send(msg)
		return;
	}

	allPharmacies, err := pharmacies.GetAllPharmacies()
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Some error happen: " + err.Error())
		bot.Send(msg)
		return;
	}

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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, seatemp_commands_message)
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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, pollution_commands_message)
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