package pollutionservice

import (
	"github.com/Ahineya/cyprushelper/helpers/storage"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
	"github.com/Ahineya/cyprushelper/dataproviders/pollution"
	"github.com/Ahineya/cyprushelper/helpers/logger"
	"errors"
)

const (
	update_time = "2000"
)

func Start(bot *tgbotapi.BotAPI) {
	pollutionServiceChannel := make(chan string)
	createPollutionService(pollutionServiceChannel)
	go func() {
		for pollutionData := range pollutionServiceChannel {

			// Get all chat Ids for limassol
			chats, err := storage.GetPollutionSubscribersForCity("limassol")
			if err != nil {
				continue
			}
			// Send messages to all subscribers, ~3 messages in a second
			for _, chatId := range chats {
				msg := tgbotapi.NewMessage(chatId, pollutionData)
				bot.Send(msg)
				time.Sleep(time.Millisecond * 300)
			}

			// TODO: Set a ticker for batch sending updates to ids

		}
	}()
}

func createPollutionService(syncChan chan string) {
	logger.Info("PollutionService", "Initialized")
	ticker := time.NewTicker(time.Second)
	cachedPollutionLevel := "High"
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
						if needToSendPollutionUpdate(cachedPollutionLevel, pollutionLevel) {
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

func needToSendPollutionUpdate(oldPollutionLevel string, newPollutionLevel string) bool {
	if (oldPollutionLevel != newPollutionLevel) {
		if oldPollutionLevel == "Low" || oldPollutionLevel == "Moderate" {
			return newPollutionLevel == "High" || newPollutionLevel == "Very high"
		} else if oldPollutionLevel == "High" || oldPollutionLevel == "Very high" {
			return newPollutionLevel == "Low" || newPollutionLevel == "Moderate"
		}
	}

	return false
}

func ManageSubscription(chatId int64, city string, action string) (bool, error) {

	if action == "subscribe" {
		err := storage.SubscribeToPollution(chatId, city)
		if err != nil {
			return false, err
		}
		return true, nil
	} else if action == "unsubscribe" {
		err := storage.UnsubscribeFromPollution(chatId, city)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, errors.New("Incorrect subscription command")
}