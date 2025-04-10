package handlers

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lealre/gator/internal/commands"
	"github.com/lealre/gator/internal/database"
	"github.com/lib/pq"
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
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}

}

func scrapeFeeds(s *commands.State) error {
	fmt.Println("Starting scrapeFeeds...") // Log entry into function

	ctx := context.Background()
	feedId, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("error getting feed to fetch: %w", err)
	}

	fmt.Println("Feed ID:", feedId)

	feed, err := s.Db.MarkFeedFetched(ctx, feedId)
	if err != nil {
		return fmt.Errorf("error marking feed as fetched: %w", err)
	}

	fmt.Println("Feed fetched successfully")

	feedData, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching URL %s: %w", feed.Url, err)
	}

	fmt.Println("Fetched feed data successfully")

	for _, item := range feedData.Channel.Item {
		fmt.Printf("Processing item: %s\n", item.Title)

		publishedDate, err := parseRSSDate(item.PubDate)
		if err != nil {
			return fmt.Errorf("error parsing the post date %s: %w", item.PubDate, err)
		}

		newPost := database.CreatePostParams{
			ID:          uuid.New(),
			Title:       sql.NullString{String: item.Title, Valid: true},
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			FeedID:      feedId,
			PublishedAt: publishedDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		post, err := s.Db.CreatePost(ctx, newPost)
		if err != nil {
			// Check for unique constraint violation error (23505)
			if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
				fmt.Printf("Skipping URL %s. Already exists.\n", item.Link)
				continue
			}
			return fmt.Errorf("error storing post for link %s: %w", item.Link, err)
		}

		fmt.Printf("- Link %s stored in database\n", post.Url)
	}

	return nil
}

func parseRSSDate(dateStr string) (time.Time, error) {
	layout := "Mon, 02 Jan 2006 15:04:05 -0700"
	return time.Parse(layout, dateStr)
}
