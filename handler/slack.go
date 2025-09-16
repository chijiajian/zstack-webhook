package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/chijiajian/zstack-webhook/config"
)

// SlackMessage defines the structure for a Slack message payload.
type SlackMessage struct {
	Text string `json:"text"`
}

// SendToSlack sends a message to a Slack Webhook URL.
func SendToSlack(message string, cfg config.WebHookConfig) {
	if cfg.URL == "" {
		log.Println("Slack URL is missing in the configuration, skipping message.")
		return
	}

	slackMessage := SlackMessage{
		Text: message,
	}

	jsonBody, err := json.Marshal(slackMessage)
	if err != nil {
		log.Println("Failed to encode Slack message to JSON format:", err)
		return
	}

	resp, err := http.Post(cfg.URL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("Failed to send Slack message:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Slack Webhook responded with an error status code: %d", resp.StatusCode)
		return
	}
	log.Println("Message successfully sent to Slack channel!")
}
