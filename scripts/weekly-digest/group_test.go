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

	assert.Equal(t, "2023", got["2023-16"].Year)
	assert.Equal(t, "16", got["2023-16"].WeekNumber)
	assert.Len(t, got["2023-16"].Items, 1)
	assert.Equal(t, main.Timestamp{time.Unix(1682257852, 0)}, got["2023-16"].Items[0].TimeUpdated)

	assert.Equal(t, "2023", got["2023-17"].Year)
	assert.Equal(t, "17", got["2023-17"].WeekNumber)
	assert.Len(t, got["2023-17"].Items, 1)
	assert.Equal(t, main.Timestamp{time.Unix(1682805248, 0)}, got["2023-17"].Items[0].TimeUpdated)
}
