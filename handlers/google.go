package handlers

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/yansigit/cmd-gpt/lib"
	"google.golang.org/api/option"
)

func HandleGoogle(prompt string) (string, error) {
	ctx := context.Background()
	cfg, err := lib.LoadConfig()
	if err != nil {
		return "", err
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GoogleAPIKey))
	if err != nil {
		fmt.Println("Error occurred during Google API client creation:", err)
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		fmt.Println("Error occurred during Google API call:", err)
		return "", err
	}
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(text), nil
		}
	}
	return "", fmt.Errorf("no response from google")
}
