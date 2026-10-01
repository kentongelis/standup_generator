package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/kentongelis/standup/collector"
	"github.com/kentongelis/standup/config"
	"github.com/spf13/cobra"
)

// Root command, just runs the standup report
var rootCmd = &cobra.Command{
	Use:   "standup",
	Short: "Generate your daily standup",
	RunE:  runStandup,
}

// runStandup loads config, collects activity, and prints it
func runStandup(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	var c collector.Collector = collector.NewGit(cfg.Repos, cfg.GitEmail)

	since := time.Now().Add(-24 * time.Hour) // TODO: last-workday logic
	activities, err := c.Collect(cmd.Context(), since)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning (%s): %v\n", c.Name(), err) // don't fail the whole run over a collector error
	}

	for _, a := range activities {
		fmt.Printf("[%s] %s: %s\n", c.Name(), a.Kind, a.Title)
	}
	return nil
}

// Execute runs the CLI and exits non-zero on failure
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
