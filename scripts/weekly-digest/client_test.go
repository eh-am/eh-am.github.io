package main_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	main "weekly-digest"
	wd "weekly-digest"

	"github.com/stretchr/testify/assert"
)

func loadFixture(filepath string) []byte {
	data, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return data
}

func TestClient(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fixture := loadFixture("testdata/example-response.json")
		w.Write(fixture)
	}))
	defer svr.Close()

	c := wd.NewClient(svr.URL, "", "")
	pr, err := c.Do(time.Now())
	if err != nil {
		t.Errorf("expected err to be nil but got '%v'", err)
	}

	if len(pr.List) != 1 {
		t.Errorf("expected list to have 1 item but found '%d'", len(pr.List))
	}

	assert.Equal(t, pr.List["3856231477"], wd.Item{
		Id:          "3856231477",
		GivenUrl:    "https://twitter.com/thockin/status/1652112019485773824",
		GivenTitle:  "Tim Hockin (thockin.yaml) no Twitter: \"Lessons about API design that I inte",
		TimeUpdated: main.Timestamp{time.Unix(1682805248, 0)},
	})
}
