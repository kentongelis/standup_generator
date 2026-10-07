package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/kentongelis/standup/collector"
)

// Format builds the standup text from a list of activities
func Format(activities []collector.Activity, now time.Time) string {
	// Copy so we don't reorder the caller's slice, then sort oldest to newest
	sorted := make([]collector.Activity, len(activities))
	copy(sorted, activities)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})

	// Anything before midnight today counts as yesterday
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var yesterday, today []collector.Activity
	for _, a := range sorted {
		if a.Timestamp.Before(midnight) {
			yesterday = append(yesterday, a)
		} else {
			today = append(today, a)
		}
	}

	var b strings.Builder // efficient way to build up a long string
	writeSection(&b, "Yesterday", yesterday)
	writeSection(&b, "Today", today)
	return b.String()
}

// writeSection writes one titled section, grouped by repo
func writeSection(b *strings.Builder, title string, acts []collector.Activity) {
	b.WriteString(title)
	b.WriteString("\n")

	if len(acts) == 0 {
		b.WriteString(" • Nothing recorded\n\n")
		return
	}

	var repos []string                          // repo names in the order they first appear
	byRepo := map[string][]collector.Activity{} // activities for each repo
	for _, a := range acts {
		if _, seen := byRepo[a.Repo]; !seen {
			repos = append(repos, a.Repo) // remember order, since
		}
		byRepo[a.Repo] = append(byRepo[a.Repo], a)
	}

	for _, repo := range repos {
		fmt.Fprintf(b, " %s\n", repo)
		for _, a := range byRepo[repo] {
			fmt.Fprintf(b, "     • %s\n", describe(a))
		}
	}
	b.WriteString("\n")
}

// describe turns an activity into a readable line
func describe(a collector.Activity) string {
	switch a.Kind {
	case collector.KindCommit:
		return "Committed: " + a.Title
	case collector.KindPROpened:
		return "Opened PR " + a.Title
	case collector.KindPRMerged:
		return "Merged PR " + a.Title
	case collector.KindPRReviewed:
		return "Reviewed PR " + a.Title
	case collector.KindIssueClosed:
		return "Closed issue " + a.Title
	default:
		return a.Title // unknown kinds still show up
	}
}
