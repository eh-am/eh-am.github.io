package internal_test

import (
	"testing"
	"time"
	"weekly_digest/internal"
)

// TODO: write a table test
func TestWeekFinder(t *testing.T) {
	// Monday
	monday := time.Date(2022, 12, 5, 0, 0, 0, 0, time.UTC)

	start, end := internal.FindWeek(monday)

	_, _, startDay := start.Date()
	if startDay != 5 {
		t.Fatal("start is not monday, 5th: ", start)
	}

	_, _, endDay := end.Date()
	if endDay != 11 {
		t.Fatal("end is not sunday, 11: ", end, end.Weekday())
	}

	// Sunday
	sunday := time.Date(2022, 12, 4, 0, 0, 0, 0, time.UTC)
	start, end = internal.FindWeek(sunday)
	_, _, startDay = start.Date()
	if startDay != 28 {
		t.Fatal("start is not monday, 28th: ", start)
	}

	_, _, endDay = end.Date()
	if endDay != 4 {
		t.Fatal("end is not sunday, 4: ", end, end.Weekday())
	}
}
