/*
Copyright © 2024 SEONGBIN YOON <yoonsb@outlook.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	cnst "github.com/yansigit/cmd-gpt/constants"
	"github.com/yansigit/cmd-gpt/handlers"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate shell code",
	Long:  `This command generates shell code based on the given prompt.`,
	Run: func(cmd *cobra.Command, args []string) {
		prompt, _ := cmd.Flags().GetString("prompt")
		shell, _ := cmd.Flags().GetString("shell")

		shellType, err := cnst.ToShellType(shell)
		if err != nil {
			panic(err)
		}

		if err := handlers.HandleChat(cnst.ShellCodeGen, shellType, prompt, ""); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringP("prompt", "p", "", "Prompt to generate content")
	genCmd.Flags().StringP("shell", "s", "", "Shell you are using (bash|powershell)")
}
