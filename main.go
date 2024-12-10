package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			playerIds := []string{
				"212587672",
				"233689177",
				"234052235",
				"209541269",
			}

			maps := map[string]interface{}{}

			for _, playerId := range playerIds {
				requestPlayerAcc(playerId)
				resp := requestGiftCode(playerId, update.Message.Text)
				maps[playerId] = resp
			}

			j, _ := json.MarshalIndent(maps, " ", "  ")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%+v\n", string(j)))
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func requestPlayerAcc(playerId string) string {
	const secret = "tB87#kPtkxqOS2"
	var currentTime = time.Now().UTC().UnixMilli()

	param := url.Values{}
	param.Set("fid", playerId)
	param.Set("time", fmt.Sprintf("%d", currentTime))
	param.Set("sign", GetMD5Hash(fmt.Sprintf("fid=%s&time=%d%s", playerId, currentTime, secret)))

	// fmt.Printf("ASD %s")

	payload := bytes.NewBufferString(param.Encode())

	request, err := http.NewRequest(http.MethodPost, "https://wos-giftcode-api.centurygame.com/api/player", payload)
	if err != nil {
		// handle error
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		// handle error
	}

	defer response.Body.Close()

	var data interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		// handle error
	}

	return fmt.Sprintf("%+v", data)
}

func requestGiftCode(playerId string, giftCode string) string {
	const secret = "tB87#kPtkxqOS2"
	var currentTime = time.Now().UTC().UnixMilli()

	param := url.Values{}
	param.Set("fid", playerId)
	param.Set("time", fmt.Sprintf("%d", currentTime))
	param.Set("cdk", giftCode)
	param.Set("sign", GetMD5Hash(fmt.Sprintf("cdk=%s&fid=%s&time=%d%s", giftCode, playerId, currentTime, secret)))

	// fmt.Printf("ASD %s")

	payload := bytes.NewBufferString(param.Encode())

	request, err := http.NewRequest(http.MethodPost, "https://wos-giftcode-api.centurygame.com/api/gift_code", payload)
	if err != nil {
		// handle error
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		// handle error
	}

	defer response.Body.Close()

	var data interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		// handle error
	}

	return fmt.Sprintf("%+v", data)
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
