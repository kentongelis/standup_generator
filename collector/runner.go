package collector

import (
	"context"
	"sync"
	"time"
)

// Result holds what one collector returned
type Result struct {
	Name       string     // which collector this came from
	Activities []Activity // what it found
	Err        error      // any error it reported
}

// RunAll runs every collector at the same time and returns one Result per collector
func RunAll(ctx context.Context, collectors []Collector, since time.Time) []Result {
	results := make([]Result, len(collectors)) // one slot per collector
	var wg sync.WaitGroup

	for i, c := range collectors {
		wg.Add(1) // one more goroutine to wait for

		go func(i int, c Collector) {
			defer wg.Done() // mark this goroutine finished, even if it panics or returns early

			acts, err := c.Collect(ctx, since)
			results[i] = Result{Name: c.Name(), Activities: acts, Err: err} // write only to our own slot
		}(i, c)
	}

	wg.Wait() // block until every goroutine has called Done\
	return results
}
