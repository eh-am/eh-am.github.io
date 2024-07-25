package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"golang.org/x/exp/slices"
)

var ErrGroupMismatch = errors.New("groups mismatch")

// Write writes each to a different file
func Write(dest string, groups map[string]GroupedItems) error {
	os.MkdirAll(dest, os.ModePerm)

	for _, g := range groups {

		// Pad since we gonna sort using the file system
		// TODO: check if file already exists and only add new items?
		padded := fmt.Sprintf("%s-%02s", g.Year, g.WeekNumber)
		filename := filepath.Join(dest, padded+".json")

		// File already exists?
		maybeExistingData, err := loadFile(filename)
		if err != nil {
			return err
		}

		out, err := merge(maybeExistingData, g)
		if err != nil {
			return err
		}

		file, err := json.MarshalIndent(out, "", " ")
		if err != nil {
			return err
		}

		fmt.Println("writing", filename)
		err = ioutil.WriteFile(filename, file, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func loadFile(filename string) (GroupedItems, error) {
	var out GroupedItems

	dat, err := os.ReadFile(filename)
	if errors.Is(err, fs.ErrNotExist) {
		return out, nil
	}

	err = json.Unmarshal(dat, &out)
	if err != nil {
		return out, err
		//		log.Printf("error unmarshalling file: '%s'. still continuing\n", err.Error())
	}

	return out, nil
}

func merge(m1, m2 GroupedItems) (GroupedItems, error) {
	var out GroupedItems

	// m2 is always defined, m1 not because it may come from a new file
	if m1.SortKey == "" {
		m1.SortKey = m2.SortKey
	}
	if m1.WeekNumber == "" {
		m1.WeekNumber = m2.WeekNumber
	}
	if m1.Year == "" {
		m1.Year = m2.Year
	}

	// Some sanity checks
	if m1.SortKey != m2.SortKey {
		return out, fmt.Errorf("%s: sortKeyMismatch: %s, %s", ErrGroupMismatch, m1.SortKey, m2.SortKey)
	}

	if m1.WeekNumber != m2.WeekNumber {
		return out, fmt.Errorf("%s: week number: %d, %d", ErrGroupMismatch, m1.WeekNumber, m2.WeekNumber)
	}

	if m1.Year != m2.Year {
		return out, fmt.Errorf("%s: week number: %d, %d", ErrGroupMismatch, m1.Year, m2.Year)
	}

	out.SortKey = m1.SortKey
	out.WeekNumber = m1.WeekNumber
	out.Year = m1.Year

	out.Items = append(m1.Items, m2.Items...)

	// Remove duplicates
	// Equivalent to unix's | sort | uniq
	// This makes it definitely slower
	sort.Slice(out.Items, func(i, j int) bool {
		// TODO: maybe we should do something smart like time updated?
		return out.Items[i].TimeAdded.Time.Before(out.Items[j].TimeAdded.Time)
	})

	out.Items = slices.CompactFunc(out.Items, func(el1, el2 Item) bool {
		return el1.Id == el2.Id
	})

	//	out := make(map[string]GroupedItems)
	return out, nil
}
