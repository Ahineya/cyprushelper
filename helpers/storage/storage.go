package storage

import (
	"os"
	"github.com/orchestrate-io/gorc"
	"github.com/Ahineya/cyprushelper/helpers/utils"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/Ahineya/cyprushelper/helpers/logger"
	"github.com/Ahineya/cyprushelper/dataproviders/pollution"
	"fmt"
)

type Chats struct {
	Ids Ids `json:"ids"`
}

type Ids []int64

var storage_token string
var cachedChatIds Ids
var chats_collection_name string

func UpdateChats(message *tgbotapi.Message) {
	chats, err := GetChatIds()
	if err != nil {
		logger.Error("STORAGE", err.Error())
	}

	if !utils.Int64InSlice(message.Chat.ID, chats) {
		chats = append(chats, message.Chat.ID)

		c := gorc.NewClient(storage_token)
		_, err = c.Put("stats", chats_collection_name, Chats{chats})
		if err != nil {
			logger.Error("STORAGE", err.Error())
			return
		}

		logger.Info("STORAGE", "Chats updated")
	} else {
		logger.Info("STORAGE", "No new chats")
	}
}

func GetChatIds() ([]int64, error) {
	if len(cachedChatIds) == 0 {
		if chats_collection_name == "" {
			env := os.Getenv("ENV")
			if env == "PROD" {
				chats_collection_name = "chats"
			} else {
				chats_collection_name = "chats-test"
			}
		}

		if storage_token == "" {
			storage_token = os.Getenv("STORAGE_TOKEN")
		}

		c := gorc.NewClient(storage_token)
		result, err := c.Get("stats", chats_collection_name)
		if err != nil {
			logger.Error("STORAGE", err.Error())
			return []int64{}, err
		}

		var chats Chats
		result.Value(&chats)

		cachedChatIds = chats.Ids

		return chats.Ids, nil
	}

	return cachedChatIds, nil
}

func SubscribeToPollution(chatId int64, city string) error {
	city = pollution.NormalizeCity(city)
	if storage_token == "" {
		storage_token = os.Getenv("STORAGE_TOKEN")
	}
	c := gorc.NewClient(storage_token)
	result, err := c.Get("pollution-subscriptions", city)
	if err != nil {
		logger.Error("STORAGE", "Subscribe to pollution error: " + err.Error())
		return err
	}

	var chats Chats
	result.Value(&chats)

	logger.Log("STP", fmt.Sprintf("%s", chats.Ids))

	found := false
	for _, chId := range chats.Ids {
		if chId == chatId {
			found = true
			break
		}
	}

	if !found {
		chats.Ids = append(chats.Ids, chatId)

		_, err = c.Put("pollution-subscriptions", city, chats)
		if err != nil {
			logger.Error("STORAGE", "Subscribe to pollution error: " + err.Error())
			return err
		}
	}

	return nil
}

func UnsubscribeFromPollution(chatId int64, city string) error {
	city = pollution.NormalizeCity(city)
	if storage_token == "" {
		storage_token = os.Getenv("STORAGE_TOKEN")
	}
	c := gorc.NewClient(storage_token)
	result, err := c.Get("pollution-subscriptions", city)
	if err != nil {
		logger.Error("STORAGE", "Unsubscribe from pollution error: " + err.Error())
		return err
	}

	var chats Chats
	result.Value(&chats)

	logger.Log("STP", fmt.Sprintf("%s", chats.Ids))

	found := false
	var idx int
	for index, chId := range chats.Ids {
		if chId == chatId {
			found = true
			idx = index
			break
		}
	}

	if found {
		chats.Ids = append(chats.Ids[:idx], chats.Ids[idx+1:]...)

		_, err = c.Put("pollution-subscriptions", city, chats)
		if err != nil {
			logger.Error("STORAGE", "Unsubscribe from pollution error: " + err.Error())
			return err
		}
	}

	return nil
}

func GetPollutionSubscribersForCity(city string) (Ids, error) {
	city = pollution.NormalizeCity(city)
	if storage_token == "" {
		storage_token = os.Getenv("STORAGE_TOKEN")
	}
	c := gorc.NewClient(storage_token)
	result, err := c.Get("pollution-subscriptions", city)
	if err != nil {
		logger.Error("STORAGE", "Unsubscribe from pollution error: " + err.Error())
		return Ids{}, err
	}

	var chats Chats
	result.Value(&chats)
	return chats.Ids, nil
}