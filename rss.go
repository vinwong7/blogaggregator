package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"log"
	"time"

	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/vinwong7/blogaggregator/internal/database"
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

		var nullTitle bool
		var nullDescription bool
		var nullPublishedAt bool

		if item.Title == "" {
			nullTitle = false
		} else {
			nullTitle = true
		}
		if item.Description == "" {
			nullDescription = false
		} else {
			nullDescription = true
		}
		if item.PubDate == "" {
			nullPublishedAt = false
		} else {
			nullPublishedAt = true
		}

		parsedDate, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			fmt.Printf("Time Parse failed. Here is the date: %v\n", item.PubDate)
			continue
		}

		itemTitle := sql.NullString{String: item.Title, Valid: nullTitle}
		itemDescription := sql.NullString{String: item.Description, Valid: nullDescription}
		itemPublishedAt := sql.NullTime{Time: parsedDate, Valid: nullPublishedAt}

		_, err = s.db.CreatePost(context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       itemTitle,
				Url:         item.Link,
				Description: itemDescription,
				PublishedAt: itemPublishedAt,
				FeedID:      nextFeed.ID,
			})
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("%v: %v\n", rssFeed.Channel.Title, item.Title)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {

	time_between_reqs := 30 * time.Second

	ticker := time.NewTicker(time_between_reqs)

	fmt.Println("Collecting feeds every 30 secs...")

	for ; ; <-ticker.C {
		scrapeFeeds(s, cmd)
		fmt.Println("Looping...")

	}

}
