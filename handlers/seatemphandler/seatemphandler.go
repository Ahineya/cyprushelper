package seatemphandler

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"github.com/Ahineya/cyprushelper/dataproviders/seatemp"
)

const (
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
)

func Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update, tokens []string) {

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

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sea temperature in " + strings.Join(tokens[1:], " ") + " is " + seatemp)
	bot.Send(msg)

}