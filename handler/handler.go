package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/chijiajian/zstack-webhook/config"
)

type InboundWebhookPayload struct {
	Sections []map[string]interface{} `json:"sections"`
}

func WebhookHandler(cfg *config.Config, outputFormat string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		var payload InboundWebhookPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Printf("Failed to parse request body: %v", err)
			http.Error(w, "Invalid request body format", http.StatusBadRequest)
			return
		}

		for _, section := range payload.Sections {
			go handleSingleAlert(section, cfg.Webhooks, outputFormat)
		}

		w.WriteHeader(http.StatusOK)
	}
}

/*
func handleSingleAlert(alertData map[string]interface{}, webhooks []config.WebhookTarget, outputFormat string) {

	var alertText string

	switch strings.ToLower(outputFormat) {
	case "json":
		jsonBytes, err := json.MarshalIndent(alertData, "", "  ")
		if err != nil {
			log.Printf("Failed to encode alert data to JSON: %v", err)
			return
		}
		alertText = fmt.Sprintf("New Alert!\nRaw JSON：\n%s", string(jsonBytes))

	default:
		var details []string

		var keys []string
		for key := range alertData {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := alertData[key]
			details = append(details, fmt.Sprintf("%s: %v", key, value))
		}

		alertDetailsText := strings.Join(details, "\n")
		alertText = fmt.Sprintf("New Alert!\n%s", alertDetailsText)
	}

	log.Printf("Preparing to send message: \n%s", alertText)

	for _, webhookTarget := range webhooks {
		go func(target config.WebhookTarget) {
			switch target.Type {
			case "slack":
				SendToSlack(alertText, target.Config)
			case "telegram":
				SendToTelegram(alertText, target.Config)
			default:
				log.Printf("Unknown webhook type: %s", target.Type)
			}
		}(webhookTarget)
	}
}
*/

// handleSingleAlert processes and sends a single alert.
func handleSingleAlert(alertData map[string]interface{}, webhooks []config.WebhookTarget, outputFormat string) {
	log.Printf("Preparing to send message from alert data:\n%+v", alertData)

	// Iterate over all configured webhook targets and send messages concurrently.
	for _, webhookTarget := range webhooks {
		go func(target config.WebhookTarget) {
			var filteredData map[string]interface{}
			var alertText string

			if len(target.Config.Fields) > 0 {
				filteredData = getFilteredAlertData(alertData, target.Config.Fields)
			} else {
				filteredData = alertData
			}

			switch strings.ToLower(outputFormat) {
			case "json":
				jsonBytes, err := json.MarshalIndent(filteredData, "", "  ")
				if err != nil {
					log.Printf("Failed to encode filtered alert data to JSON for %s webhook: %v", target.Type, err)
					return
				}
				alertText = fmt.Sprintf("📢 New Alert!\nRaw JSON:\n%s", string(jsonBytes))

			default: // Defaults to "text".
				var details []string

				var keys []string
				for key := range filteredData {
					keys = append(keys, key)
				}
				sort.Strings(keys)

				for _, key := range keys {
					value := filteredData[key]
					details = append(details, fmt.Sprintf("%s: %v", key, value))
				}

				alertDetailsText := strings.Join(details, "\n")
				alertText = fmt.Sprintf("📢 New Alert!\n%s", alertDetailsText)
			}

			log.Printf("Sending to %s webhook. Message content:\n%s", target.Type, alertText)

			switch target.Type {
			case "slack":
				SendToSlack(alertText, target.Config)
			case "telegram":
				SendToTelegram(alertText, target.Config)
			case "dingtalk":
				SendToDingTalk(alertText, target.Config)
			default:
				log.Printf("Unknown webhook type: %s", target.Type)
			}
		}(webhookTarget)
	}
}

// getFilteredAlertData filters the alert data based on the provided fields list.
func getFilteredAlertData(originalData map[string]interface{}, fields []string) map[string]interface{} {
	filteredData := make(map[string]interface{})
	for _, field := range fields {
		if value, ok := originalData[field]; ok {
			filteredData[field] = value
		}
	}
	return filteredData
}
