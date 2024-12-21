package handlers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/yansigit/cmd-gpt/config"
	cnst "github.com/yansigit/cmd-gpt/constants"
)

func HandleChat(commandType cnst.CommandType, shell cnst.ShellType, input string, provider string) error {
	godotenv.Load()
	var res string
	prompt := input
	if commandType == cnst.ShellCodeGen {
		if shell == cnst.None {
			return fmt.Errorf("shell param is empty")
		}
		prompt = fmt.Sprintf("Generate a %s command for: %s. Only output the command, no explanations.", shell, input)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	if provider == "" {
		provider = cfg.DefaultProvider
	}

	switch provider {
	case cnst.OpenAI:
		res, err = HandleOpenAI(prompt)
		if err != nil {
			return err
		}
	case cnst.Anthropic:
		res, err = HandleAnthropic(prompt)
		if err != nil {
			return err
		}
	case cnst.Google:
		res, err = HandleGoogle(prompt)
		if err != nil {
			return err
		}
	case cnst.OpenRouter:
		res, err = HandleOpenRouter(prompt)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid provider: %s", provider)
	}

	if shell != cnst.None {
		err := executeCommand(res, shell)
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
