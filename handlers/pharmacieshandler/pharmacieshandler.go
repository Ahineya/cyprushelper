package pharmacieshandler

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/Ahineya/cyprushelper/dataproviders/pharmacies"
	"strings"
)

const (
	pharmacies_commands_message = `
	Please, define city: /pharmacies city_name
	Or use commands for available cities:

	Limassol: /pharmacies_limassol
	Paphos: /pharmacies_pafos
	Larnaka: /pharmacies_larnaka
	Nicosia: /pharmacies_nicosia
	Famagusta: /pharmacies_famagusta
	`
)

func Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update, tokens []string) {

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
