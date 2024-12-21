package handlers

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/yansigit/cmd-gpt/lib"
)

func HandleOpenAI(systemPrompt string, prompt string) (string, error) {
	cfg, err := lib.LoadConfig()
	if err != nil {
		return "", err
	}

	client := openai.NewClient(cfg.OpenAIKey)
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})

	if err != nil {
		fmt.Println("Error occurred during OpenAI API call:", err)
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
