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

type pollutionServiceData struct {
	City string
	Message string
}

func Start(bot *tgbotapi.BotAPI) {
	pollutionServiceChannel := make(chan pollutionServiceData)
	createPollutionService(pollutionServiceChannel)
	go func() {
		for pollutionData := range pollutionServiceChannel {

			// Get all chat Ids for limassol
			chats, err := storage.GetPollutionSubscribersForCity(pollutionData.City)
			if err != nil {
				continue
			}
			// TODO: Create a sending queue, because we have events for four cities here ?
			for _, chatId := range chats {
				msg := tgbotapi.NewMessage(chatId, pollutionData.Message)
				bot.Send(msg)
				time.Sleep(time.Millisecond * 300)
			}

		}
	}()
}

func createPollutionService(syncChan chan pollutionServiceData) {
	logger.Info("PollutionService", "Initialized")
	ticker := time.NewTicker(time.Second)

	pollutionLevelsCache := make(map[string]string)
	pollutionLevelsCache["LIMTRA"] = "High"
	pollutionLevelsCache["LARTRA"] = "High"
	pollutionLevelsCache["PAFTRA"] = "High"
	pollutionLevelsCache["NICTRA"] = "High"

	go func() {
		for t := range ticker.C {
			if t.Format("0405") == update_time {
				logger.Info("PollutionService", "Getting updates")
				getPollutionUpdates(syncChan, "limassol", pollutionLevelsCache)
				getPollutionUpdates(syncChan, "larnaka", pollutionLevelsCache)
				getPollutionUpdates(syncChan, "paphos", pollutionLevelsCache)
				getPollutionUpdates(syncChan, "nicosia", pollutionLevelsCache)
			}
		}
	}()
}

func getPollutionUpdates(syncChan chan pollutionServiceData, city string, pollutionLevelsCache map[string]string) {
	pollutionResult, err := pollution.GetPollution(city)
	if err != nil {
		logger.Warn("PollutionService", city + ": " + err.Error())
		return
	}

	normalizedCity := pollution.NormalizeCity(city)

	for _, pollutionData := range pollutionResult.Data {
		if pollutionData.PollutantCode == "PM10" {
			pollutionLevel := pollution.GetPollutionLevel(pollutionData.PollutantCode, pollutionData.Value)
			if needToSendPollutionUpdate(pollutionLevelsCache[normalizedCity], pollutionLevel) {
				pollutionLevelsCache[normalizedCity] = pollutionLevel
				logger.Info("PollutionService", "Pollution level for " + city +" changed, sending updates")

				syncChan <- pollutionServiceData{city, "Current pollution level in " + city + " changed to " + pollutionLevel}
			} else {
				logger.Info("PollutionService", "Pollution level for " + city +" not changed")
			}
		}
	}
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