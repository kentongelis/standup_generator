package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/kentongelis/standup/collector"
	"github.com/kentongelis/standup/config"
	"github.com/kentongelis/standup/workday"

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
	cfg, err := config.Load() // read ~/.standup.yaml and env vars
	if err != nil {
		return err
	}

	since := workday.LastWorkday((time.Now())) // last-workday logic

	// Every source of activity, all treated the same way
	collectors := []collector.Collector{
		collector.NewGit(cfg.Repos, cfg.GitEmail),
		collector.NewGitHub(cfg.GitHubToken, cfg.GitHubUsername),
	}

	results := collector.RunAll(cmd.Context(), collectors, since) // run them all at once

	// Merge results and report any failures
	var activities []collector.Activity
	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "warning (%s): %v\n", r.Name, r.Err) // warn but keep going
		}
		activities = append(activities, r.Activities...) // safe: back on a single go routine
	}

	fmt.Printf("Activity since %s\n\n", since.Format("Mon Jan 2"))
	for _, a := range activities {
		fmt.Printf("[%s] %s: %s\n", a.Repo, a.Kind, a.Title)
	}
	return nil
}

// Execute runs the CLI and exits non-zero on failure
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
