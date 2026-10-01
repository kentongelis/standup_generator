package collector

import (
	"context"
	"time"
)

// Kind type describes what type of work an Activity represents
type Kind string

const (
	KindCommit      Kind = "commit"
	KindPROpened    Kind = "pr_opened"
	KindPRMerged    Kind = "pr_merged"
	KindPRReviewed  Kind = "pr_reviewed"
	KindIssueCLosed Kind = "issue_closed"
)

// The normalized piece of work
type Activity struct {
	Source    string    // which collector
	Kind      Kind      // what happened
	Title     string    // commit message, pr tittle, issue title
	URL       string    // link to it (empty for local commits)
	Repo      string    // repo name
	Timestamp time.Time // when it happened
}

// Collector is anything that can report a user's activity since a given time
type Collector interface {
	Name() string
	Collect(ctx context.Context, since time.Time) ([]Activity, error)
}
