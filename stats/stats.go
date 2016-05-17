package stats

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"fmt"
	"os"
	"github.com/botanio/sdk/go"
	"strings"
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
		fmt.Printf("Asynchonous: %+v\n", ans)
		ch <- true
	})

	<-ch
}
