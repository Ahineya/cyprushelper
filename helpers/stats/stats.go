package stats

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"github.com/botanio/sdk/go"
	"strings"
	"github.com/Ahineya/cyprushelper/helpers/logger"
)

// TODO: Make proper metrics
type Message struct {
	SomeMetric    int
	AnotherMetric int
}

var botan_token string

func Track(message *tgbotapi.Message) {
	ch := make(chan bool)

	if (botan_token == "") {
		botan_token = os.Getenv("BOTAN_TOKEN")
	}
	bot := botan.New(botan_token)

	tokens := strings.Fields(message.Text)
	command := tokens[0]

	bot.TrackAsync(message.From.ID, Message{100, 500}, command, func(ans botan.Answer, err []error) {
		if len(err) == 0 {
			logger.Info("BOTAN", ans.Status + " " + ans.Info)
		} else {
			logger.Warn("BOTAN", ans.Status + " " + ans.Info)
		}
		ch <- true
	})

	<-ch
}