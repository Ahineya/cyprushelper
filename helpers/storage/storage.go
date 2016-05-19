package storage

import (
	"os"
	"github.com/orchestrate-io/gorc"
	"fmt"
	"github.com/Ahineya/cyprushelper/helpers/utils"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Chats struct {
	Ids []int64 `json:"ids"`
}

var storage_token string
var cachedChats []int64
var chats_collection_name string
func UpdateChats(message *tgbotapi.Message) {
	chats, err := GetChatIds()
	if err != nil {
		fmt.Printf(err.Error())
	}

	if !utils.Int64InSlice(message.Chat.ID, chats) {
		chats = append(chats, message.Chat.ID)

		c := gorc.NewClient(storage_token)
		_, err = c.Put("stats", chats_collection_name, Chats{chats})
		if err != nil {
			fmt.Println("[STORAGE]:" + err.Error())
			return
		}

		fmt.Println("[STORAGE]: Chats updated");
	} else {
		fmt.Println("[STORAGE]: No new chats");
	}
}

func GetChatIds() ([]int64, error) {
	if len(cachedChats) == 0 {
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
			fmt.Println("[STORAGE]:" + err.Error())
			return []int64{}, err
		}

		var chats Chats
		result.Value(&chats)

		cachedChats = chats.Ids

		return chats.Ids, nil
	}

	return cachedChats, nil
}