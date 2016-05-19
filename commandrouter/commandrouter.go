package commandrouter

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/Ahineya/cyprushelper/helpers/stats"
	"github.com/Ahineya/cyprushelper/helpers/storage"
	"strings"
	"github.com/Ahineya/cyprushelper/handlers/pollutionhandler"
	"github.com/Ahineya/cyprushelper/handlers/seatemphandler"
	"github.com/Ahineya/cyprushelper/handlers/pharmacieshandler"
)

var Messages = map[string]string{
	"Start": "Hi! I am Cyprus Helper bot!",
	"Help": `You can use the following commands:
	/pharmacies
	/seatemp
	/pollution
	`,
}

const (
	// TODO: define as environment variable
	bot_name = "CyprusHelper_bot"
)

func Route(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	//if os.Getenv("ENV") == "PROD" {
	stats.Track(update.Message)
	go storage.UpdateChats(update.Message)
	//}

	tokens := strings.Fields(update.Message.Text)
	command := tokens[0]

	// Checking for composite commands like /pharmacies_limassol instead of /pharmacies limassol
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
		pharmacieshandler.Handle(bot, update, tokens)
	case "/pharmacies" + bot_name:
		pharmacieshandler.Handle(bot, update, tokens)
	case "/seatemp":
		seatemphandler.Handle(bot, update, tokens)
	case "/seatemp" + bot_name:
		seatemphandler.Handle(bot, update, tokens)
	case "/pollution":
		pollutionhandler.Handle(bot, update, tokens)
	case "/pollution" + bot_name:
		pollutionhandler.Handle(bot, update, tokens)
	}
}

