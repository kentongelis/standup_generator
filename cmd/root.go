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
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	since := workday.LastWorkday((time.Now())) // last-workday logic

	// Every source of activity, all treated the sanem way
	collectors := []collector.Collector{
		collector.NewGit(cfg.Repos, cfg.GitEmail),
		collector.NewGitHub(cfg.GitHubToken, cfg.GitHubUsername),
	}

	var activities []collector.Activity
	for _, c := range collectors {
		acts, err := c.Collect(cmd.Context(), since)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning (%s): %v\n", c.Name(), err) // warn but keep going
		}
		activities = append(activities, acts...) // keep whatever this collector found
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
