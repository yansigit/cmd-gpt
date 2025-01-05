package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/charmbracelet/huh"

	cnst "github.com/yansigit/cmd-gpt/constants"
	"github.com/yansigit/cmd-gpt/lib"
)

var logger *lib.Logger
var cfg *lib.Config

func HandleChat(commandType cnst.CommandType, shell string, input string, provider string) error {
	systemPrompt := "You are a terminal assistant. Main purpose is to help the user to execute commands in the terminal. You will be given a command and you will execute it in the terminal. You will only output the full working, and valid commands without any explanations. No markdown output. Just plain commands."
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
		platform := runtime.GOOS
		if platform == "darwin" {
			platform = "macos"
		}
		prompt = fmt.Sprintf("I am on %s. Generate %s shell code for: %s.", runtime.GOOS, shell, input)
	}

	if provider == "" {
		provider = cfg.DefaultProvider
	}

	res, err := generateCommand(provider, systemPrompt, prompt)
	if err != nil {
		return err
	}

	err = executeCommandWithConfirmation(res, shell, "Execute the command?", provider, systemPrompt, prompt)
	if err != nil {
		return err
	}

	return nil
}

func generateCommand(provider string, systemPrompt string, prompt string) (string, error) {
	var res string
	var err error
	switch provider {
	case cnst.OpenAI:
		res, err = HandleOpenAI(systemPrompt, prompt)
	case cnst.Anthropic:
		res, err = HandleAnthropic(systemPrompt, prompt)
	case cnst.Google:
		res, err = HandleGoogle(systemPrompt, prompt)
	case cnst.OpenRouter:
		res, err = HandleOpenRouter(systemPrompt, prompt)
	default:
		return "", fmt.Errorf("invalid provider: %s", provider)
	}
	return trimMarkdownCodeTags(res), err
}

func executeCommandWithConfirmation(command string, shell string, message string, provider string, systemPrompt string, prompt string) error {
	logger.Success("Command from GPT: ", command)
	var confirmation bool
	err := huh.NewConfirm().Title(message).Value(&confirmation).Run()
	if err != nil {
		logger.Errorf("error receiving confirmation: %v", err)
		return err
	}
	if confirmation {
		if err := executeCommand(command, shell); err != nil {
			logger.Error("Error occurred during executing the command:", err)
			return err
		}
	} else {
		logger.Error("Skipped executing the command by the user.")
		var regenerate bool
		err = huh.NewConfirm().Title("Regenerate the command?").Value(&regenerate).Run()
		if err != nil {
			return fmt.Errorf("error receiving confirmation for regeneration: %v", err)
		}
		if regenerate {
			var followUpPrompt string
			huh.NewInput().
				Title("Follow up prompt").
				Prompt(">").
				Value(&followUpPrompt).Run()
			prompt = fmt.Sprintf("%s\n%s", prompt, followUpPrompt)
			logger.Info("Regenerating the command...")
			res, err := generateCommand(provider, systemPrompt, prompt)
			if err != nil {
				return err
			}
			err = executeCommandWithConfirmation(res, shell, "Execute the new command?", provider, systemPrompt, prompt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func trimMarkdownCodeTags(input string) string {
	openTagRegex := regexp.MustCompile("(?m)^```.*[ \n]")
	result := openTagRegex.ReplaceAllString(input, "")

	closeTagRegex := regexp.MustCompile("(?m)^```$")
	result = closeTagRegex.ReplaceAllString(result, "")

	return strings.TrimSpace(result)
}

func executeCommand(command string, shell string) error {
	// Split command into lines to handle multiline commands
	commands := strings.Split(command, "\n")

	// If it's a multiline command, create a temporary shell file and execute that
	if len(commands) > 1 {
		f, err := os.CreateTemp("", "*."+strings.Split(shell, "/")[len(strings.Split(shell, "/"))-1])
		if err != nil {
			return fmt.Errorf("failed to create temporary file: %v", err)
		}
		defer os.Remove(f.Name()) // clean up

		_, err = f.WriteString(command)
		if err != nil {
			return fmt.Errorf("failed to write to temporary file: %v", err)
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("failed to close temporary file: %v", err)
		}

		cmd := exec.Command(shell, f.Name())
		cmd.Dir = "." // Set working directory to current directory
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to execute command '%s': %v", command, err)
		}
	} else {
		for _, cmdStr := range commands {
			cmdStr = strings.TrimSpace(cmdStr)
			if cmdStr == "" {
				continue // Skip empty lines
			}

			var cmd *exec.Cmd
			if shell == "powershell" {
				cmd = exec.Command("powershell", "-Command", cmdStr)
			} else {
				cmd = exec.Command(shell, "-c", cmdStr)
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to execute command '%s': %v", cmdStr, err)
			}
		}
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
