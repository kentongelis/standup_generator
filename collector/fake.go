package collector

import (
	"context"
	"time"
)

type Fake struct{}

// Don't need to implement Collector interface, just need to the two methods with matching signatures

func (Fake) Name() string {
	return "fake"
}

func (Fake) Collect(ctx context.Context, since time.Time) ([]Activity, error) {
	return []Activity{
		{Source: "fake", Kind: KindCommit, Title: "Set up Cobra", Repo: "standup", Timestamp: time.Now()},
		{Source: "fake", Kind: KindPRMerged, Title: "Add config loading", Repo: "standup", Timestamp: time.Now()},
	}, nil
}
