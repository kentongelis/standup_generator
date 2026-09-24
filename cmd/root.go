package cmd

import (
	"fmt"
	"os"

	"github.com/kentongelis/standup/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "standup",
	Short: "Generate your daily standup",
	RunE:  runStandup,
}

func runStandup(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	fmt.Println("user:", cfg.GitHubUsername)
	fmt.Println("repos:", cfg.Repos)
	// fmt.Println("token set:", cfg.GitHubToken != "")
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
