package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/yansigit/cmd-gpt/constants"
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
	logger = GetLogger()
}

func InitConfig() error {
	logger.Info("Initializing config...")
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		logger.Info("Config file not found, creating a new one...")
		SetProvider()
	} else {
		logger.Info("Config file already exists. Skipping initialization...")
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

func SetProvider() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	var provider string
	huh.NewSelect[string]().Title("Select a provider").Options(
		huh.NewOption("OpenRouter", constants.OpenRouter),
		huh.NewOption("OpenAI", constants.OpenAI),
		huh.NewOption("Anthropic", constants.Anthropic),
		huh.NewOption("Google", constants.Google),
	).Value(&provider).Run()

	var apiKey string
	huh.NewInput().Title(fmt.Sprintf("enter your %s api key", provider)).EchoMode(huh.EchoModePassword).Value(&apiKey).Run()

	// Check if the API key is valid
	if apiKey == "" {
		return fmt.Errorf("API key is empty")
	}

	switch provider {
	case constants.OpenAI:
		cfg.OpenAIKey = apiKey
	case constants.Anthropic:
		cfg.AnthropicKey = apiKey
	case constants.Google:
		cfg.GoogleAPIKey = apiKey
	case constants.OpenRouter:
		cfg.OpenRouterKey = apiKey
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}

	cfg.DefaultProvider = provider
	return SaveConfig(cfg)
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
