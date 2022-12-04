package internal

import (
	"time"
)

func daysSinceMonday(curr time.Time) int {
	// For example, remember that time.Sunday is 0
	// So 0 + 6 days passed since last Monday
	return int(curr.Weekday()+6) % 7
}

// FindWeek finds the beginning (monday) and end (sunday) of a week
func FindWeek(d time.Time) (start time.Time, end time.Time) {
	monday := d.Add(-24 * time.Hour * time.Duration(daysSinceMonday(d)))
	sunday := monday.Add(24 * time.Hour * 6)

	return monday, sunday
}
