package config

import (
	"os"
)

type Config struct {
	OpenAIKey    string
	AnthropicKey string
}

func LoadConfig() Config {
	return Config{
		OpenAIKey:    os.Getenv("OPENAI_API_KEY"),
		AnthropicKey: os.Getenv("ANTHROPIC_API_KEY"),
	}
}
