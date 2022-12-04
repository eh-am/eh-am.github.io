package internal_test

import (
	"os"
	"testing"
	"time"
	wd "weekly_digest/internal"
)

type MockArticlesGetter struct {
	items []wd.PocketItem
	err   error
}

func (m *MockArticlesGetter) GetArticles(from time.Time, until time.Time) ([]wd.PocketItem, error) {
	return m.items, m.err
}

func TestWeeklyDigest(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	mag := MockArticlesGetter{}
	//	digest := wd.NewArticlesDigest(tmpDir, mag)
	//	digest.Run(from time.Time, until time.Time)
}
