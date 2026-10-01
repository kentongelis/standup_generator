package workday

import (
	"testing"
	"time"
)

func TestLastWorkday(t *testing.T) {
	// each case: a fake "now" and the start time expected back
	tests := []struct {
		name string
		now  time.Time
		want time.Time
	}{
		{"monday goes back to friday", date(2026, 10, 5, 9), date(2026, 10, 2, 0)},
		{"tuesday goes back to monday", date(2026, 10, 6, 9), date(2026, 10, 5, 0)},
		{"saturday goes back to friday", date(2026, 10, 3, 9), date(2026, 10, 2, 0)},
		{"sunday goes back to friday", date(2026, 10, 4, 9), date(2026, 10, 2, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LastWorkday(tt.now)
			if !got.Equal(tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

// helper functino for building a time at the start of an hour
func date(y int, m time.Month, d, h int) time.Time {
	return time.Date(y, m, d, h, 0, 0, 0, time.Local)
}
