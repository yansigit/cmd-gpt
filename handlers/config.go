package handlers

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/yansigit/cmd-gpt/config"
	"github.com/yansigit/cmd-gpt/constants"
)

func HandleSetProvider() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	var provider string
	huh.NewSelect[string]().Title("Select a provider").Options(
		huh.NewOption("OpenAI", constants.OpenAI),
		huh.NewOption("Anthropic", constants.Anthropic),
		huh.NewOption("Google", constants.Google),
		huh.NewOption("OpenRouter", constants.OpenRouter),
	).Value(&provider)

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
	return config.SaveConfig(cfg)
}
