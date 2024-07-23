package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Write writes each to a different file
func Write(dest string, groups map[string]GroupedItems) error {
	os.MkdirAll(dest, os.ModePerm)

	for _, g := range groups {
		file, err := json.MarshalIndent(g, "", " ")
		if err != nil {
			return err
		}

		// Pad since we gonna sort using the file system
		// TODO: check if file already exists and only add new items?
		padded := fmt.Sprintf("%s-%02s", g.Year, g.WeekNumber)
		err = ioutil.WriteFile(filepath.Join(dest, padded+".json"), file, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
