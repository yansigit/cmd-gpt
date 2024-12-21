package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yansigit/cmd-gpt/lib"
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

func HandleOpenRouter(systemPrompt string, prompt string) (string, error) {
	cfg, err := lib.LoadConfig()
	if err != nil {
		return "", err
	}

	if cfg.DefaultModel == "" {
		cfg.DefaultModel = "gemini-2.0-flash"
	}

	apiKey := cfg.OpenRouterKey
	if apiKey == "" {
		return "", logger.Errorf("OPENROUTER api key is not set")
	}

	requestBody, err := json.Marshal(OpenRouterRequest{
		Model: cfg.DefaultModel,
		Messages: []OpenRouterMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
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
		return "", logger.Errorf("openrouter api error: %s", resp.Status)
	}

	var openRouterResponse OpenRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&openRouterResponse); err != nil {
		return "", err
	}

	if len(openRouterResponse.Choices) > 0 {
		return openRouterResponse.Choices[0].Message.Content, nil
	}

	return "", logger.Errorf("no response from openrouter")
}
