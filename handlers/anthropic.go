package handlers

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/joho/godotenv"
)

func HandleAnthropic(prompt string) (string, error) {
	godotenv.Load()
	client := anthropic.NewClient()
	msg, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5HaikuLatest),
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewAssistantMessage(anthropic.NewTextBlock("You are a terminal assistant. Main purpose is to help the user to execute commands in the terminal. You will be given a command and you will execute it in the terminal. You will only output the command, no explanations.")),
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
	})

	if err != nil {
		fmt.Println("Error occurred during Anthropic API call:", err)
		return "", err
	}
	return msg.Content[0].Text, nil
}
