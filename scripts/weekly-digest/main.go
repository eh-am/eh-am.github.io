package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/markusmobius/go-dateparser"
)

func main() {
	outPtr := flag.String("out", "./output", "the output dir")
	flag.Parse()

	args := flag.Args()
	if len(args) <= 0 {
		log.Fatalf("missing mandatory argument")
	}

	var since time.Time
	if args[0] == "all" {
		since = time.Time{}
	} else {
		args := strings.Join(args, " ")
		dt, err := dateparser.Parse(nil, args)
		if err != nil {
			log.Fatalf("unknown arg: '%s'", args)
		}
		since = dt.Time
	}

	err := run(since, *outPtr)
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

func run(from time.Time, outDir string) error {
	fmt.Println("running from", from)
	url := "https://getpocket.com/v3/get"
	consumerKey, accessToken := getAccess()

	client := NewClient(url, consumerKey, accessToken)
	clientRes, err := client.Do(from)
	if err != nil {
		return fmt.Errorf("error calling client: %v", err)
	}

	// clientRes
	fmt.Printf("Found '%d' items\n", len(clientRes.List))

	grouped := GroupByISOWeek(clientRes.List)
	err = Write(outDir, grouped)
	if err != nil {
		return fmt.Errorf("error writing to dir: %v", err)
	}

	return nil
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
