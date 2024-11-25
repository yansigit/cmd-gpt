package handlers

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/yansigit/cmd-gpt/config"
	cnst "github.com/yansigit/cmd-gpt/constants"
)

func HandleChat(commandType cnst.CommandType, shell cnst.ShellType, input string) error {
	_config := config.LoadConfig()
	anthropicApiKey := _config.AnthropicKey

	client := anthropic.NewClient(
		option.WithAPIKey(anthropicApiKey),
	)

	prompt := input
	if commandType == cnst.ShellCodeGen {
		if shell == cnst.None {
			return fmt.Errorf("shell param is empty")
		}
		prompt = fmt.Sprintf("Generate a %s command for: %s. Only output the command, no explanations.", shell, input)
	}

	msg, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5HaikuLatest),
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewAssistantMessage(anthropic.NewTextBlock("You are a terminal assistant. Main purpose is to help the user to execute commands in the terminal. You will be given a command and you will execute it in the terminal. You will only output the command, no explanations.")),
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
	})

	if err != nil {
		fmt.Println("Error occurred during API call:", err)
		return err
	}

	// prettyJSON, _ := json.MarshalIndent(msg.Content, "", "  ")
	// fmt.Println(string(prettyJSON))

	res := msg.Content[0].Text

	if shell != cnst.None {
		err = executeCommand(res, shell)
		if err != nil {
			fmt.Println("Error occurred during executing the command:", err)
			return err
		}
	} else {
		fmt.Println("Response: ", res)
	}

	return nil
}

func executeCommand(command string, shell cnst.ShellType) error {
	cmd := exec.Command(string(cnst.Bash), "-c", command)
	if shell == cnst.PowerShell {
		cmd = exec.Command(string(cnst.PowerShell), "-Command", command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}

	return nil
}
