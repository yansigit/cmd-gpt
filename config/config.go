package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

type Config struct {
	OpenAIKey       string `json:"openai_key"`
	AnthropicKey    string `json:"anthropic_key"`
	OpenRouterKey   string `json:"openrouter_key"`
	GoogleAPIKey    string `json:"google_api_key"`
	DefaultProvider string `json:"default_provider"`
	DefaultModel    string `json:"default_model"`
}

var (
	configDir  string
	configPath string
)

func init() {
	var err error
	configDir, err = os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	configPath = filepath.Join(configDir, "cmd-gpt", "config.json")
}

func InitConfig() error {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		var apiKey string
		huh.NewInput().Title("enter your openrouter api key").EchoMode(huh.EchoModePassword).Value(&apiKey).Run()
		cfg := &Config{
			DefaultProvider: "openrouter",
			DefaultModel:    "gemini-2.0-flash",
			OpenRouterKey:   apiKey,
		}
		if err := SaveConfig(cfg); err != nil {
			return err
		}
	}
	return nil
}

func LoadConfig() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(configDir, "cmd-gpt", "config.json")

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func SaveConfig(cfg *Config) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(configDir, "cmd-gpt", "config.json")

	err = os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(cfg); err != nil {
		return err
	}

	return nil
}

func (c *Config) String() string {
	out, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(out)
}
