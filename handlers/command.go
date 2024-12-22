package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/huh"

	cnst "github.com/yansigit/cmd-gpt/constants"
	"github.com/yansigit/cmd-gpt/lib"
)

var logger *lib.Logger
var cfg *lib.Config

func HandleChat(commandType cnst.CommandType, shell string, input string, provider string) error {
	var res string
	systemPrompt := "You are a terminal assistant. Main purpose is to help the user to execute commands in the terminal. You will be given a command and you will execute it in the terminal. You will only output the command, no explanations."
	prompt := input
	if commandType == cnst.ShellCodeGen {
		if shell == "" {
			logger.Warn("Shell not specified. Will detect shell automatically.")
			shell = os.Getenv("SHELL")
			if shell == "" {
				if os.PathSeparator == '\\' {
					shell = "powershell" // Default to 'powershell' for Windows
					logger.Warn("Detected Windows OS. Defaulting to 'Powershell'.")
				} else {
					shell = "sh" // Default to 'sh' if no shell is detected
					logger.Warn("Unable to detect shell. Defaulting to 'sh'.")
				}
			} else {
				shell = filepath.Base(shell)
				logger.Info("Detected shell: ", shell)
			}

		}
		prompt = fmt.Sprintf("Generate a %s command for: %s.", shell, input)
	}

	if provider == "" {
		provider = cfg.DefaultProvider
	}

	var err error
	switch provider {
	case cnst.OpenAI:
		res, err = HandleOpenAI(systemPrompt, prompt)
		if err != nil {
			return err
		}
	case cnst.Anthropic:
		res, err = HandleAnthropic(systemPrompt, prompt)
		if err != nil {
			return err
		}
	case cnst.Google:
		res, err = HandleGoogle(systemPrompt, prompt)
		if err != nil {
			return err
		}
	case cnst.OpenRouter:
		res, err = HandleOpenRouter(systemPrompt, prompt)
		if err != nil {
			return err
		}
	default:
		return logger.Errorf("invalid provider: %s", provider)
	}

	logger.Success("Command from GPT: ", res)
	var confirmation bool
	err = huh.NewConfirm().Title("Execute the command?").Value(&confirmation).Run()
	if err != nil {
		logger.Errorf("error receiving confirmation: %v", err)
		return err
	}
	if confirmation {
		if err := executeCommand(res, shell); err != nil {
			logger.Error("Error occurred during executing the command:", err)
			return err
		}
	} else {
		logger.Error("Skipped executing the command by the user.")
	}

	return nil
}

func executeCommand(command string, shell string) error {
	var cmd *exec.Cmd
	if shell == "powershell" {
		cmd = exec.Command("powershell", "-Command", command)
	} else {
		cmd = exec.Command(shell, "-c", command)
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

	var err error
	cfg, err = lib.LoadConfig()
	if err != nil {
		logger.Fatal("Error loading config:", err)
	}
}
