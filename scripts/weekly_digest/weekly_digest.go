package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	fmt.Println("vim-go")
}

func Run() {
	var from time.Time
	var until time.Time

	// Get from pocket API items in the range (from, until)
	_ = from
	_ = until

	// Create the file if not exists

	// Update the file with the new items,
	// Not touching existing changes
}

type PocketItem struct {
	Id    string
	URL   string
	Title string
}

type ArticlesGetter struct{}

func (ag *ArticlesGetter) GetArticles(from time.Time, until time.Time) ([]PocketItem, error) {
	return []PocketItem{}, nil
}

type ArticlesDigest struct {
	pathPrefix string
	ag         ArticlesGetter
}

func NewArticlesDigest(pathPrefix string, ag ArticlesGetter) *ArticlesDigest {
	return &ArticlesDigest{
		pathPrefix: pathPrefix,
		ag:         ag,
	}
}

func (ad *ArticlesDigest) Run(from time.Time, until time.Time) error {
	// TODO: check if from is a monday and until is a sunday

	// Get from pocket API items in the range (from, until)
	allItems, err := ad.ag.GetArticles(from, until)
	if err != nil {
		return err
	}

	targetFilename := ad.targetFilename(from, until)
	file, err := os.Open(targetFilename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	_ = allItems

	// Find new items
	newItems, err := ad.findNewItems(scanner, allItems)
	if err != nil {
		return err
	}

	for _, item := range newItems {
		n, err := file.WriteString(fmt.Sprintf("[%s](%s) <!-- auto:id:%s -->\n", item.Title, item.URL, item.Id))
		if n <= 0 {
			return fmt.Errorf("written %d characters which is invalid", n)
		}
		if err != nil {
			return err
		}
	}
	// Update the file with the new items,
	// Not touching existing changes
	return nil
}

func (ad *ArticlesDigest) targetFilename(from time.Time, until time.Time) string {
	return filepath.Join(ad.pathPrefix, fmt.Sprintf("from-%s-to-%s.md", from.Format("2006-01-02"), until.Format("2006-01-02")))
}

// findNewItems removes in-place items that are already in the document
func (ad *ArticlesDigest) findNewItems(scanner *bufio.Scanner, newItems []PocketItem) ([]PocketItem, error) {
	foundIds := make(map[string]struct{})

	// Pass once, finding all the existing items
	for scanner.Scan() {
		// We are assuming there are no long lines (> 64 KB)
		// https://cs.opensource.google/go/go/+/refs/tags/go1.19.3:src/bufio/scan.go;l=82
		if strings.Contains(scanner.Text(), "auto:") {
			// Items
			foundIds["myid"] = struct{}{}
		}
	}

	brandNewItems := make([]PocketItem, len(newItems)-len(foundIds))

	i := 0
	for _, item := range newItems {
		if _, ok := foundIds[item.Id]; ok {
			// That items is already there, so let's skip it
			continue
		}

		brandNewItems[i] = item
		i = i + 1
	}

	return brandNewItems, nil
}
