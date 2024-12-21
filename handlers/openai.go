package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func HandleOpenAI(prompt string) (string, error) {
	godotenv.Load()
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a terminal assistant. Main purpose is to help the user to execute commands in the terminal. You will be given a command and you will execute it in the terminal. You will only output the command, no explanations.",
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
