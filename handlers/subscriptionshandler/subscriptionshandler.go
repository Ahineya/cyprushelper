package subscriptionshandler

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func Handle(bot *tgbotapi.BotAPI, update tgbotapi.Update, tokens []string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	var keyboard [][]tgbotapi.KeyboardButton
	for _, city := range []string{"limassol", "larnaka", "paphos", "nicosia"} {
		keyboard = append(keyboard, []tgbotapi.KeyboardButton{
			tgbotapi.KeyboardButton{Text: "/pollution " + city + " subscribe"},
		}, )
		keyboard = append(keyboard, []tgbotapi.KeyboardButton{
			tgbotapi.KeyboardButton{Text: "/pollution " + city + " unsubscribe"},
		}, )
	}

	msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
		Keyboard: keyboard,
		ResizeKeyboard: true,
		OneTimeKeyboard: true,
	}
	bot.Send(msg)
}