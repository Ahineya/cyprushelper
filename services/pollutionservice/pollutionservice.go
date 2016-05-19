package pollutionservice

import (
	"github.com/Ahineya/cyprushelper/helpers/storage"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
	"github.com/Ahineya/cyprushelper/dataproviders/pollution"
	"github.com/Ahineya/cyprushelper/helpers/logger"
)

const (
	update_time = "2000"
)

func Start(bot *tgbotapi.BotAPI) {
	pollutionServiceChannel := make(chan string)
	createPollutionService(pollutionServiceChannel)
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

func createPollutionService(syncChan chan string) {
	logger.Info("PollutionService", "Initialized")
	ticker := time.NewTicker(time.Second)
	cachedPollutionLevel := ""
	go func() {
		for t := range ticker.C {
			if t.Format("0405") == update_time {
				logger.Info("PollutionService", "Getting updates")
				pollutionResult, err := pollution.GetPollution("limassol")
				if err != nil {
					logger.Warn("PollutionService", err.Error())
				}

				for _, pollutionData := range pollutionResult.Data {
					if pollutionData.PollutantCode == "PM10" {
						pollutionLevel := pollution.GetPollutionLevel(pollutionData.PollutantCode, pollutionData.Value)
						if pollutionLevel != cachedPollutionLevel {
							cachedPollutionLevel = pollutionLevel
							logger.Info("PollutionService", "Pollution level changed, sending updates")
							syncChan <- "Current pollution level in Limassol changed to " + pollutionLevel
						}
					}
				}
			}
		}
	}()
}