package handler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/chijiajian/zstack-webhook/config"
)

func SendToTelegram(message string, cfg config.WebHookConfig) {
	if cfg.BotToken == "" || cfg.ChatID == "" {
		log.Println("Telegram configuration is incomplete, skipping.")
		return
	}

	// Construct the Telegram API URL.
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.BotToken)

	// Use net/url.Values to build URL-encoded form data.
	formData := url.Values{
		"chat_id": {cfg.ChatID},
		"text":    {message},
	}

	// Create an HTTP client.
	client := &http.Client{}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		log.Printf("Failed to create Telegram request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send Telegram message: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Successfully sent to Telegram.")

	if resp.StatusCode != http.StatusOK {
		log.Printf("Telegram API returned an error status code: %d", resp.StatusCode)
		return
	}

	log.Println("Message successfully sent to Telegram!")
}

// TelegramMessage struct is not needed with this method
type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}
