package internal_test

import (
	"reflect"
	"testing"
	"time"
	"weekly_digest/internal"
)

// TODO: write a table test
func TestWeekFinder(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		// Previous week
		{input: "2022-12-04", want: []string{"2022-11-28", "2022-12-04"}},

		// New week
		{input: "2022-12-05", want: []string{"2022-12-05", "2022-12-11"}},
		{input: "2022-12-06", want: []string{"2022-12-05", "2022-12-11"}},
		{input: "2022-12-07", want: []string{"2022-12-05", "2022-12-11"}},
		{input: "2022-12-08", want: []string{"2022-12-05", "2022-12-11"}},
		{input: "2022-12-09", want: []string{"2022-12-05", "2022-12-11"}},
		{input: "2022-12-10", want: []string{"2022-12-05", "2022-12-11"}},
		{input: "2022-12-11", want: []string{"2022-12-05", "2022-12-11"}},

		// New week
		{input: "2022-12-12", want: []string{"2022-12-12", "2022-12-18"}},
	}

	for _, tc := range tests {
		d, err := time.Parse("2006-01-02", tc.input)
		if err != nil {
			panic(err)
		}
		start, end := internal.FindWeek(d)

		gotStart := start.Format("2006-01-02")
		gotEnd := end.Format("2006-01-02")

		if !reflect.DeepEqual(tc.want[0], gotStart) {
			t.Fatalf("expected: %v, got: %v", tc.want[0], gotStart)
		}

		if !reflect.DeepEqual(tc.want[1], gotEnd) {
			t.Fatalf("expected: %v, got: %v", tc.want[1], gotEnd)
		}
	}
}
