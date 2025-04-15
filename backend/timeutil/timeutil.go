package timeutil

import "time"

// Normalize a time to UTC and round it to microseconds (to match PostgreSQL).
func Normalize(t time.Time) time.Time {
	return t.UTC().Round(time.Microsecond)
}

// Return the current time (normalized to UTC and rounded to microseconds).
func Now() time.Time {
	return Normalize(time.Now())
}
