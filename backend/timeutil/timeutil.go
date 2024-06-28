package timeutil

import "time"

// Return the current UTC time rounded to microseconds (to match PostgreSQL).
func Now() time.Time {
	return time.Now().UTC().Round(time.Microsecond)
}
