package main_test

import (
	"testing"
	"time"
	main "weekly-digest"
	wd "weekly-digest"

	"github.com/stretchr/testify/assert"
)

func TestGroup(t *testing.T) {
	got := wd.GroupByISOWeek(map[string]wd.Item{
		"x": wd.Item{
			TimeUpdated: main.Timestamp{time.Unix(1682805248, 0)},
		},
		"y": wd.Item{
			TimeUpdated: main.Timestamp{time.Unix(1682257852, 0)},
		},
	})

	want := map[string]wd.GroupedItems{
		"2023-16": wd.GroupedItems{
			Year:       "2023",
			WeekNumber: "16",
			Items: []wd.Item{
				{
					TimeUpdated: main.Timestamp{time.Unix(1682257852, 0)},
				},
			},
		},
		"2023-17": wd.GroupedItems{
			Year:       "2023",
			WeekNumber: "17",
			Items: []wd.Item{
				{
					TimeUpdated: main.Timestamp{time.Unix(1682805248, 0)},
				},
			},
		},
	}

	assert.Equal(t, want, got)
}
