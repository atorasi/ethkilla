package utils

import (
	"bytes"
	"encoding/json"
	"ethkilla/src/constants"
	"fmt"
	"log"

	http "github.com/Danny-Dasilva/fhttp"
)

func SendTelegramMessage(format string, args ...interface{}) {
	if !constants.SETTINGS.TelegramAlerts {
		return
	}

	url := "https://api.telegram.org/bot" + constants.SETTINGS.BotToken + "/sendMessage"

	payload, _ := json.Marshal(&telegramMessage{
		ChatID:                constants.SETTINGS.ChatID,
		Text:                  fmt.Sprintf(format, args...),
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
	})
	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err == nil && response.StatusCode >= 200 && response.StatusCode < 300 {
		return
	}
	defer response.Body.Close()
	log.Println("An error or bad response code in sending Telegram message")
}
