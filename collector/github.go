package collector

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v74/github"
)

// fails to compile if GitHub doesn't match the Collector interface
var _ Collector = (*GitHub)(nil)

// GitHub collects PRs, reviews, and issues from GitHub
type GitHub struct {
	client   *github.Client // authenticated API client
	username string         // github username
}

// NewGitHub creates a GitHub collector
func NewGitHub(token, username string) *GitHub {
	return &GitHub{
		client:   github.NewClient(nil).WithAuthToken(token),
		username: username,
	}
}

// Name identifies this collector in warnings
func (g *GitHub) Name() string { return "github" }

// Collect gathers your GitHub activity since the given time
func (g *GitHub) Collect(ctx context.Context, since time.Time) ([]Activity, error) {
	if g.username == "" {
		return nil, errors.New("github_username is not set in config")
	}

	date := since.UTC().Format(time.RFC3339) // format GitHub search understands
	var activities []Activity
	var errs []error

	// PRs opened or merged
	authored, err := g.search(ctx, fmt.Sprintf("is:pr author:%s updated:>=%s", g.username, date))
	if err != nil {
		errs = append(errs, fmt.Errorf("authored PRs: %w", err))
	}
	for _, pr := range authored {
		if a, ok := authoredPR(pr, since); ok {
			activities = append(activities, a)
		}
	}

	// PRs reviewed excluding your own
	reviewed, err := g.search(ctx, fmt.Sprintf("is:pr reviewed-by:%s -author:%s updated:>=%s", g.username, g.username, date))
	if err != nil {
		errs = append(errs, fmt.Errorf("reviewed PRs: %w", err))
	}
	for _, pr := range reviewed {
		activities = append(activities, toActivity(pr, KindPRReviewed, pr.GetUpdatedAt().Time))
	}

	// issues assigned that were closed
	closed, err := g.search(ctx, fmt.Sprintf("is:issue is:closed assignee:%s closed:>=%s", g.username, date))
	if err != nil {
		errs = append(errs, fmt.Errorf("closed issues: %w", err))
	}
	for _, issue := range closed {
		activities = append(activities, toActivity(issue, KindIssueClosed, issue.GetClosedAt().Time))
	}

	return activities, errors.Join(errs...) // nil if nothing fails
}

//  search runs one GitHub search query and returns the matching items
func (g *GitHub) search(ctx context.Context, query string) ([]*github.Issue, error) {
	opts := &github.SearchOptions{ListOptions: github.ListOptions{PerPage: 100}}
	result, _, err := g.client.Search.Issues(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	return result.Issues, nil
}

// authoredPR decides wither your PR was merged or opened in the window
func authoredPR(pr *github.Issue, since time.Time) (Activity, bool) {
	if merged := pr.GetPullRequestLinks().GetMergedAt(); merged.After(since) {
		return toActivity(pr, KindPRMerged, merged.Time), true
	}
	if created := pr.GetCreatedAt(); created.After(since) {
		return toActivity(pr, KindPROpened, created.Time), true
	}
	return Activity{}, false // only updated, not opened or merged
}

// toActivity converts a GitHub search result into an Activity
func toActivity(i *github.Issue, kind Kind, when time.Time) Activity {
	return Activity{
		Source:    "github",
		Kind:      kind,
		Title:     fmt.Sprintf("#%d %s", i.GetNumber(), i.GetTitle()),
		URL:       i.GetHTMLURL(),
		Repo:      repoName(i.GetRepositoryURL()),
		Timestamp: when,
	}
}

// repoName turns "https://api.github.com/repos/owner/repo" into "owner/repo"
func repoName(apiURL string) string {
	return strings.TrimPrefix(apiURL, "https://api.github.com/repos/")
}
