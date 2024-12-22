/*
Copyright © 2024 SEONGBIN YOON <yoonsb@outlook.com>
*/
package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	cnst "github.com/yansigit/cmd-gpt/constants"
	"github.com/yansigit/cmd-gpt/handlers"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate shell code",
	Long:  `This command generates shell code based on the given prompt.`,
	Run: func(cmd *cobra.Command, args []string) {
		prompt, _ := cmd.Flags().GetString("prompt")
		shell, _ := cmd.Flags().GetString("shell")

		if len(args) > 0 {
			prompt = strings.Join(args, " ")
		}

		if prompt == "" {
			logger.Fatal("Prompt is required")
		}

		if err := handlers.HandleChat(cnst.ShellCodeGen, shell, prompt, ""); err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringP("prompt", "p", "", "Prompt to generate content")
	genCmd.Flags().StringP("shell", "s", "", "Shell you are using (bash|powershell)")
}
