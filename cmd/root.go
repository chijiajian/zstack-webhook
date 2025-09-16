package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zstack-webhook",
	Short: "A webhook server for ZStack alerts",
	Long: `A flexible webhook server that listens for ZStack alerts
and forwards them to various platforms like Slack.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
