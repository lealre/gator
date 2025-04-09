package handlers

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/lealre/gator/internal/commands"
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

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting %s: %w", feedUrl, err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feedData RSSFeed
	err = xml.Unmarshal(body, &feedData)
	if err != nil {
		return nil, err
	}

	decodeHTMLEntities(&feedData)

	return &feedData, nil
}

func decodeHTMLEntities(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
}

func Aggregate(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: agg <time-between-requests>")
	}

	timeBetweenReqs := cmd.Args[0]
	duration, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		return fmt.Errorf("error parsing duration %s: %w", duration, err)
	}

	fmt.Println("Collecting feeds every ", duration)
	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func scrapeFeeds(s *commands.State) error {

	ctx := context.Background()
	feedId, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("error getting feed to fecth, %w", err)
	}

	feed, err := s.Db.MarkFeedFetched(ctx, feedId)
	if err != nil {
		return fmt.Errorf("error marking feed as fecthed: %w", err)
	}

	feedData, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching url %s: %w", feed.Url, err)
	}

	fmt.Printf("Channel title: %s\n", feedData.Channel.Title)
	fmt.Println("Feeds:")
	for _, item := range feedData.Channel.Item {
		fmt.Printf("  - %s\n", item.Title)
	}

	return nil
}
