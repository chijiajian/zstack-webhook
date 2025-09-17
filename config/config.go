package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port     int    `yaml:"port"`
	HTTPS    bool   `yaml:"https"`
	CertFile string `yaml:"cert_file,omitempty"`
	KeyFile  string `yaml:"key_file,omitempty"`
}

type WebHookConfig struct {
	URL      string   `yaml:"url,omitempty"`
	BotToken string   `yaml:"bot_token,omitempty" json:"botToken,omitempty"` //telegram bot token
	ChatID   string   `yaml:"chat_id,omitempty" json:"chatID,omitempty"`     //telegram chat id
	Secret   string   `yaml:"secret,omitempty"`                              // dingtalk robot secret
	Fields   []string `yaml:"fields,omitempty"`
}

type WebhookTarget struct {
	Type   string        `yaml:"type"`
	Config WebHookConfig `yaml:"config"`
}

type Config struct {
	Server   ServerConfig    `yaml:"server"`
	Webhooks []WebhookTarget `yaml:"webhooks"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	return &cfg, nil
}
