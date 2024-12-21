package handlers

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

func HandleAnthropic(systemPrompt string, prompt string) (string, error) {
	client := anthropic.NewClient()
	msg, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		// Model:     anthropic.F(anthropic.ModelClaude3_5HaikuLatest),
		Model:     anthropic.F(cfg.DefaultModel),
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewAssistantMessage(anthropic.NewTextBlock(systemPrompt)),
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
	})

	if err != nil {
		fmt.Println("Error occurred during Anthropic API call:", err)
		return "", err
	}
	return msg.Content[0].Text, nil
}
