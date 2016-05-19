package pollutionhandler

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"github.com/Ahineya/cyprushelper/dataproviders/pollution"
	"github.com/Ahineya/cyprushelper/services/pollutionservice"
)

const (
	pollution_commands_message = `
	Please, define city: /pollution city_name
	Or use commands for available cities:

	Limassol: /pollution_limassol
	Paphos: /pollution_pafos
	Larnaka: /pollution_larnaka
	Nicosia: /pollution_nicosia

	If you want to subscribe to dangerous pollution changes, use /pollution city_name subscribe
	If you want to unsubscribe from dangerous pollution changes, use /pollution city_name unsubscribe

	You will receive updates every hour if the pollution situation has changed noteably

	Currently subscription works only for Limassol.
	`
)

func Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update, tokens []string) {
	if len(tokens) < 2 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, pollution_commands_message)
		bot.Send(msg)
		return
	}

	city := strings.ToLower(tokens[1])

	if len(tokens) == 3 {
		command := strings.ToLower(tokens[2])
		res, err := pollutionservice.ManageSubscription(update.Message.Chat.ID, city, command)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
			bot.Send(msg)
			return;
		}

		var message string
		if res {
			message = "You have succesfully subscribed to pollution updates for the " + city
		} else {
			message = "You have succesfully unsubscribed from pollution updates for the " + city
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
		bot.Send(msg)
		return;
	}

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