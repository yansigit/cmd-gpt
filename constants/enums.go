package constants

import "errors"

type CommandType string
type ShellType string

const (
	JustChat     CommandType = "chat"
	ShellCodeGen CommandType = "shell_code_generation"
)
const (
	Bash       ShellType = "bash"
	PowerShell ShellType = "powershell"
	None       ShellType = ""
)

func ToShellType(s string) (ShellType, error) {
	switch s {
	case string(Bash):
		return Bash, nil
	case string(PowerShell):
		return PowerShell, nil
	case string(None):
		return None, nil
	default:
		return None, errors.New("invalid shell type")
	}
}
