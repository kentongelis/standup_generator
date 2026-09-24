package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "standup",
	Short: "Generate your daily standup",
	Run:   runStandup,
}

func runStandup(cmd *cobra.Command, args []string) {
	fmt.Println("standup is working")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
