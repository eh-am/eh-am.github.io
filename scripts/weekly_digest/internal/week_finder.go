package internal

import (
	"time"
)

func daysSinceMonday(curr time.Time) int {
	return int(curr.Weekday()+6) % 7
}

// FindWeek finds the beginning and end of a week
func FindWeek(d time.Time) (start time.Time, end time.Time) {
	monday := d.Add(-24 * time.Hour * time.Duration(daysSinceMonday(d)))
	sunday := monday.Add(24 * time.Hour * 6)

	return monday, sunday
}
