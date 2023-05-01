package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("either 'all' or 'last' is required")
	}

	args := os.Args[1]
	switch args {
	case "all":
		all()
	case "last":
		last()
	default:
		log.Fatalf("unknown arg: '%s'", args)
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

func all() {
	fmt.Println("running all")

	url := "https://getpocket.com/v3/get"
	consumerKey, accessToken := getAccess()

	client := NewClient(url, consumerKey, accessToken)
	clientRes, err := client.Do(time.Time{})
	if err != nil {
		panic(err)
	}
	grouped := GroupByISOWeek(clientRes.List)
	err = Write("output", grouped)
	if err != nil {
		panic(err)
	}
}

// TODO: figure out why it returns last week
func last() {
	fmt.Println("running last week")
	url := "https://getpocket.com/v3/get"
	consumerKey, accessToken := getAccess()

	client := NewClient(url, consumerKey, accessToken)
	clientRes, err := client.Do(findLastSunday(time.Now()))
	if err != nil {
		panic(err)
	}
	grouped := GroupByISOWeek(clientRes.List)
	err = Write("output", grouped)
	if err != nil {
		panic(err)
	}
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
