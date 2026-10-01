package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/kentongelis/standup/collector"
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
	_ = cfg // not used yet

	var c collector.Collector = collector.Fake{}

	since := time.Now().Add(-24 * time.Hour)
	activities, err := c.Collect(cmd.Context(), since)
	if err != nil {
		return err
	}

	for _, a := range activities {
		fmt.Printf("[%s] %s: %s\n", c.Name(), a.Kind, a.Title)
	}
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
