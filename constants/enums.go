package constants

type CommandType string
type ShellType string

const (
	JustChat     CommandType = "chat"
	ShellCodeGen CommandType = "shell_code_generation"
)

const (
	OpenAI     string = "openai"
	Anthropic  string = "anthropic"
	Google     string = "google"
	OpenRouter string = "openrouter"
)
