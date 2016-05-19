package pollutionhandler

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"github.com/Ahineya/cyprushelper/dataproviders/pollution"
)

const (
	pollution_commands_message = `
	Please, define city: /pollution city_name
	Or use commands for available cities:

	Limassol: /pollution_limassol
	Paphos: /pollution_pafos
	Larnaka: /pollution_larnaka
	Nicosia: /pollution_nicosia
	`
)

func Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update, tokens []string) {
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