package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"log"
	"time"

	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed

	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

func scrapeFeeds(s *state, _ command) error {
	//Fetch next feed based on earliest last_fetched_dt
	nextFeed, err := s.db.GetNextFeedtoFetch(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	//Mark feed as fetched by updating last_fetched_dt
	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		log.Fatal(err)
	}

	//Fetch feed using URL
	rssFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range (*rssFeed).Channel.Item {
		fmt.Printf("%v\n", item.Title)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {

	time_between_reqs := 30 * time.Second

	ticker := time.NewTicker(time_between_reqs)

	fmt.Println("Collecting feeds every 30 secs...")

	for ; ; <-ticker.C {
		scrapeFeeds(s, cmd)
		fmt.Println("....................")
	}

}
