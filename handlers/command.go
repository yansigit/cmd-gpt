package handlers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"

	cnst "github.com/yansigit/cmd-gpt/constants"
	"github.com/yansigit/cmd-gpt/lib"
)

var logger *lib.Logger

func HandleChat(commandType cnst.CommandType, shell cnst.ShellType, input string, provider string) error {
	var res string
	prompt := input
	if commandType == cnst.ShellCodeGen {
		if shell == cnst.None {
			return fmt.Errorf("shell param is empty")
		}
		prompt = fmt.Sprintf("Generate a %s command for: %s. Only output the command, no explanations.", shell, input)
	}

	cfg, err := lib.LoadConfig()
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
		logger.Info("Response: ", res)
		var confirmation bool
		err := huh.NewConfirm().Title("Execute the command?").Value(&confirmation).Run()
		if err != nil {
			logger.Errorf("error receiving confirmation: %v", err)
			return err
		}
		if confirmation {
			if err := executeCommand(res, shell); err != nil {
				logger.Errorf("Error occurred during executing the command:", err)
				return err
			}
		} else {
			logger.Error("Skipped executing the command by the user.")
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

func init() {
	logger = lib.GetLogger()
}
