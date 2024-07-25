package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	url         string
	consumerKey string
	accessToken string

	httpClient *http.Client
}

// TODO: use builder args
func NewClient(url string, consumerKey string, accessToken string) *Client {
	return &Client{
		url:         url,
		consumerKey: consumerKey,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: time.Minute * 3,
		},
	}
}

type Timestamp struct {
	time.Time
}

// UnmarshalJSON decodes an int64 timestamp into a time.Time object
// It also supports RFC3339
// https://stackoverflow.com/a/67017059
func (p *Timestamp) UnmarshalJSON(bytes []byte) error {
	// 1. Decode the bytes into an int64
	var raw string
	err := json.Unmarshal(bytes, &raw)

	if err != nil {
		return err
	}

	rawInt, err := strconv.ParseInt(raw, 10, 64)

	if err != nil {
		// We support RFC3339 too
		if errors.Is(err, strconv.ErrSyntax) {
			t, err := time.Parse(time.RFC3339, raw)
			if err != nil {
				return err
			}

			p.Time = t
			return nil
		}
	}

	// Let's assume we are always in UTC, so that build machine and local match
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		return err
	}

	// 2. Parse the unix timestamp
	p.Time = time.Unix(rawInt, 0).In(utc)
	return nil
}

type Item struct {
	Id           string    `json:"item_id"`
	GivenUrl     string    `json:"given_url"`
	GivenTitle   string    `json:"given_title"`
	ResolvedItem string    `json:"resolved_title"`
	ResolvedUrl  string    `json:"resolved_url"`
	TimeUpdated  Timestamp `json:"time_updated"`
	TimeAdded    Timestamp `json:"time_added"`
}
type PocketResponse struct {
	List map[string]Item `json:"list"`
}

func (c *Client) Do(since time.Time) (pr PocketResponse, err error) {
	type body struct {
		ConsumerKey string `json:"consumer_key"`
		AccessToken string `json:"access_token"`
		Since       int    `json:"since"`
		Sort        string `json:"sort"`
	}

	b := body{
		ConsumerKey: c.consumerKey,
		AccessToken: c.accessToken,
		Since:       int(since.Unix()),
		Sort:        "oldest",
	}

	marshalled, err := json.Marshal(b)
	if err != nil {
		return pr, err
	}

	// url "https://getpocket.com/v3/get"
	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(marshalled))
	if err != nil {
		return pr, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return pr, err
	}

	// TODO: other status code may happen
	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return pr, err
		}
		return pr, errors.New(fmt.Sprintf("statusCode: '%d'. body: '%s'", res.StatusCode, string(bodyBytes)))
	}

	err = json.NewDecoder(res.Body).Decode(&pr)
	if err != nil {
		return pr, err
	}

	return pr, nil
}
