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

/*
func handleSingleAlert(alertData map[string]interface{}, webhooks []config.WebhookTarget, outputFormat string) {

	//var alertText string
	// 构造一个字符串切片，用于存放所有键值对
	var details []string

	// 遍历 map 中的所有键值对
	for key, value := range alertData {
		// 将键值对格式化为 "键: 值" 的字符串，并添加到切片中
		details = append(details, fmt.Sprintf("%s: %v", key, value))
	}

	// 将所有详情拼接成一个完整消息
	alertDetailsText := strings.Join(details, "\n")

	log.Printf("解析成功！\n%s", alertDetailsText)

	// 构造发送到 Slack 和 Telegram 的纯文本消息
	alertText := fmt.Sprintf(
		"📢 收到新告警！\n%s",
		alertDetailsText,
	)

	// 遍历所有配置的 Webhook 目标，并并发发送
	for _, webhookTarget := range webhooks {
		go func(target config.WebhookTarget) {
			switch target.Type {
			case "slack":
				SendToSlack(alertText, target.Config)
			case "telegram":
				SendToTelegram(alertText, target.Config)
			default:
				log.Printf("未知 Webhook 类型: %s", target.Type)
			}
		}(webhookTarget)
	}
}
*/
