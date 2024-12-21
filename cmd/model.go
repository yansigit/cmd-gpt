/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/charmbracelet/huh"
	"github.com/yansigit/cmd-gpt/lib"

	"github.com/spf13/cobra"
)

// modelCmd represents the model command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Change default model to be used",
	Long: `This command is for changing default model to be used in the application.
For example, you can change default model to abcd/text-model-003 to use the model
as the default model for generating code.`,
	Run: func(cmd *cobra.Command, args []string) {
		huh.NewInput().Title("Enter model name").Value(&cfg.DefaultModel).Run()
		lib.SaveConfig(cfg)
	},
}

func init() {
	configCmd.AddCommand(modelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
