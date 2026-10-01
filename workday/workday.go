package workday

import "time"

// LastWorkday returns midnight at the start of the most recent weekday before now
func LastWorkday(now time.Time) time.Time {
	day := now.AddDate(0, 0, -1) // start from yesterday

	// step back past Saturday and Sunday
	for day.Weekday() == time.Saturday || day.Weekday() == time.Sunday {
		day = day.AddDate(0, 0, -1)
	}

	// reset to midnight but keeping the same time zone
	return time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
}
