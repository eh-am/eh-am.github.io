package main

import (
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("either 'all' | 'last' | 'last-2' is required")
	}

	var since time.Time
	args := os.Args[1]
	switch args {
	case "all":
		since = time.Time{}
	case "last":
		since = findLastSunday(time.Now())
	case "last-2":
		since = time.Now().Add(-time.Hour * 7 * 24)
	default:
		log.Fatalf("unknown arg: '%s'", args)
	}

	err := run(since)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
}

func getAccess() (string, string) {
	consumerKey := os.Getenv("POCKET_CONSUMER_KEY")
	if consumerKey == "" {
		panic("POCKET_CONSUMER_KEY is required")
	}

	accessKey := os.Getenv("POCKET_ACCESS_TOKEN")
	if accessKey == "" {
		panic("POCKET_ACCESS_TOKEN is required")
	}

	return consumerKey, accessKey
}

func run(from time.Time) error {
	url := "https://getpocket.com/v3/get"
	consumerKey, accessToken := getAccess()

	client := NewClient(url, consumerKey, accessToken)
	clientRes, err := client.Do(from)
	if err != nil {
		return err
	}
	grouped := GroupByISOWeek(clientRes.List)
	return Write("output", grouped)
}

func findLastSunday(referenceDate time.Time) time.Time {
	beginningOfDay := func(t time.Time) time.Time {
		return t.Truncate(24 * time.Hour)
	}

	if referenceDate.Weekday() == time.Sunday {
		return beginningOfDay(referenceDate)
	}

	daysToSunday := int(referenceDate.Weekday())
	lastSunday := referenceDate.AddDate(0, 0, -daysToSunday)

	return beginningOfDay(lastSunday)
}
