/*
Copyright © 2024 SEONGBIN YOON <yoonsb@outlook.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yansigit/cmd-gpt/lib"
)

var logger *lib.Logger
var cfg *lib.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cmd-gpt",
	Short: "execute gpt in your terminal",
	Long:  `This command executes gpt (llm) in your terminal.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	logger = lib.GetLogger()

	var err error
	cfg, err = lib.LoadConfig()
	if err != nil {
		logger.Fatal("Error loading config:", err)
	}
}
