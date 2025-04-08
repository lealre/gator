package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lealre/gator/internal/commands"
	"github.com/lealre/gator/internal/database"
)

func AddFeed(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	ctx := context.Background()

	user, err := s.Db.GetUserByName(ctx, s.Cfg.CurrentUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not foud: %w", err)
		}
	}

	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		Name:      name,
		Url:       url,
	}

	feed, err := s.Db.CreateFeed(ctx, newFeed)
	if err != nil {
		return fmt.Errorf("error adding new feed: %w", err)
	}

	// automatically follows the feed
	createFeedParams := database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	_, err = s.Db.CreateFeedFollow(ctx, createFeedParams)
	if err != nil {
		return fmt.Errorf("error creating the feed follow: %w", err)
	}

	fmt.Printf("url: %s\n", feed.Url)
	fmt.Printf("name: %s\n", feed.Name)

	return nil
}

func ListFeed(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: feeds")
	}

	ctx := context.Background()

	allFeeds, err := s.Db.ListAllFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}

	for _, feed := range allFeeds {
		fmt.Printf("Feed: %s\n", feed.FeedName)
		fmt.Printf("User: %s\n", feed.UserName)
		fmt.Println("------")
	}

	return nil
}
