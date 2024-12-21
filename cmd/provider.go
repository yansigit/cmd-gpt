/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yansigit/cmd-gpt/config"
	"github.com/yansigit/cmd-gpt/handlers"
)

// providerCmd represents the provider command
var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// provider, _ := cmd.Flags().GetString("provider")
		if err := handlers.HandleSetProvider(); err != nil {
			fmt.Println("Error setting provider:", err)
			os.Exit(1)
		}
		config, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
		}
		fmt.Println("Default provider set to:", config.DefaultProvider)
	},
}

func init() {
	configCmd.AddCommand(providerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// providerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// providerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
