package collector

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Checks if Git satisfies the Collector interface
var _ Collector = (*Git)(nil)

// Git collects commits from local repos
type Git struct {
	repos []string // paths to local repos
	email string   // local git author email
}

// NewGit creates a Git collector
func NewGit(repos []string, email string) *Git {
	return &Git{repos: repos, email: email}
}

// Name indentifies this collector in warnings
func (g *Git) Name() string { return "git" }

// Collect gathers commits from every repo since the given time
func (g *Git) Collect(ctx context.Context, since time.Time) ([]Activity, error) {
	var activities []Activity // all commits found
	var errs []error          // repos that failed

	for _, path := range g.repos {
		// stop early if the caller cancelled
		if err := ctx.Err(); err != nil {
			return activities, err // cancelled or timed out
		}

		acts, err := g.collectRepo(path, since)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", path, err)) // note which repo failed to collect
			continue                                             // don't stop if one repo is bad
		}
		activities = append(activities, acts...) // add this repo's commits to the list
	}

	return activities, errors.Join(errs...) // nil if nothing fails
}

// collectRepo gather commits from a single repo
func (g *Git) collectRepo(path string, since time.Time) ([]Activity, error) {
	repo, err := git.PlainOpen(path) // open the repo folder
	if err != nil {
		return nil, err
	}

	// Get commit history from all branches since the cutoff
	iter, err := repo.Log(&git.LogOptions{
		All:   true,
		Since: &since,
	})
	if err != nil {
		return nil, err
	}
	defer iter.Close() // clean up when the function ends

	repoName := filepath.Base(path)
	var acts []Activity

	// Runs once for each commit
	err = iter.ForEach(func(c *object.Commit) error {
		if !strings.EqualFold(c.Author.Email, g.email) {
			return nil // skip commits by other people
		}

		// Convert the commit into an Activity
		acts = append(acts, Activity{
			Source:    "git",
			Kind:      KindCommit,
			Title:     firstLine(c.Message),
			Repo:      repoName,
			Timestamp: c.Author.When,
		})
		return nil
	})
	return acts, err
}

// firstLine returns just the first line of a commit message
func firstLine(msg string) string {
	line, _, _ := strings.Cut(msg, "\n") // split at the first newline
	return strings.TrimSpace(line)       // remove extra whitespace
}
