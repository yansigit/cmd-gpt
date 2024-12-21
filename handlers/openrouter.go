package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/yansigit/cmd-gpt/config"
)

type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterRequest struct {
	Model    string              `json:"model"`
	Messages []OpenRouterMessage `json:"messages"`
}

type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func HandleOpenRouter(prompt string) (string, error) {
	godotenv.Load()
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}

	if cfg.DefaultModel == "" {
		cfg.DefaultModel = "gemini-2.0-flash"
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENROUTER_API_KEY environment variable is not set")
	}

	requestBody, err := json.Marshal(OpenRouterRequest{
		Model: cfg.DefaultModel,
		Messages: []OpenRouterMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/yansigit/cmd-gpt")
	req.Header.Set("X-Title", "cmd-gpt")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openrouter api error: %s", resp.Status)
	}

	var openRouterResponse OpenRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&openRouterResponse); err != nil {
		return "", err
	}

	if len(openRouterResponse.Choices) > 0 {
		return openRouterResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from openrouter")
}