package main

import (
	"fmt"
	"sort"
)

type GroupedItems struct {
	Year       string `json:"year"`
	WeekNumber string `json:"week"`
	// Key to be used when sorting in hugo
	SortKey string `json:"sortKey"`
	Items   []Item `json:"items"`
}

func GroupByISOWeek(items map[string]Item) map[string]GroupedItems {
	res := make(map[string]GroupedItems)

	for _, item := range items {
		year, week := item.TimeUpdated.ISOWeek()
		key := fmt.Sprintf("%d-%d", year, week)

		if _, ok := res[key]; !ok {
			res[key] = GroupedItems{
				Year:       fmt.Sprintf("%d", year),
				WeekNumber: fmt.Sprintf("%d", week),
				SortKey:    fmt.Sprintf("%d-%02d", year, week),
				Items:      make([]Item, 0),
			}
		}

		group := res[key]
		group.Items = append(group.Items, item)

		// This makes it definitely slower
		sort.Slice(group.Items, func(i, j int) bool {
			// TODO: maybe we should do something smart like time updated?
			return group.Items[i].TimeAdded.Time.Before(group.Items[j].TimeAdded.Time)
		})

		// Since we get a copy back, let's write it back to the map
		res[key] = group
	}

	return res
}
