// handler/dingtalk.go

package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/chijiajian/zstack-webhook/config"
)

type DingTalkMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func SendToDingTalk(message string, cfg config.WebHookConfig) {
	if cfg.URL == "" {
		log.Println("DingTalk URL is missing in the configuration, skipping message.")
		return
	}

	dingTalkMessage := DingTalkMessage{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{Content: message},
	}

	jsonBody, err := json.Marshal(dingTalkMessage)
	if err != nil {
		log.Printf("Failed to encode DingTalk message to JSON format: %v", err)
		return
	}

	requestURL := cfg.URL
	if cfg.Secret != "" {
		timestamp := time.Now().UnixNano() / 1e6
		stringToSign := fmt.Sprintf("%d\n%s", timestamp, cfg.Secret)

		h := hmac.New(sha256.New, []byte(cfg.Secret))
		h.Write([]byte(stringToSign))
		signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

		requestURL = fmt.Sprintf("%s&timestamp=%d&sign=%s", cfg.URL, timestamp, url.QueryEscape(signature))
	}

	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Failed to send DingTalk message: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("DingTalk Webhook responded with an error status code: %d", resp.StatusCode)
		return
	}

	log.Println("Message successfully sent to DingTalk!")
}
