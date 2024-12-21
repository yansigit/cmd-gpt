/*
Copyright © 2024 SEONGBIN YOON <yoonsb@outlook.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yansigit/cmd-gpt/lib"
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
		logger.Info("Setting default provider...")

		if err := lib.SetProvider(); err != nil {
			logger.Fatal("Error setting provider:", err)
			// os.Exit(1)
		}
		config, err := lib.LoadConfig()
		if err != nil {
			logger.Fatal("Error loading config:", err)
		}
		logger.Success("Default provider set to:", config.DefaultProvider)
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
