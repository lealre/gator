package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lealre/gator/internal/commands"
	"github.com/lealre/gator/internal/database"
)

func Follow(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: follow <url>")
	}

	feedUrl := cmd.Args[0]

	ctx := context.Background()

	// get feed id
	feed, err := s.Db.GetFeed(ctx, feedUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("feed '%s' not found", feedUrl)
		}
		return err
	}

	// get user id
	currentUser := s.Cfg.CurrentUser
	user, err := s.Db.GetUserByName(ctx, currentUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user '%s' not found", currentUser)
		}
		return err
	}

	createFeedParams := database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	feedFollow, err := s.Db.CreateFeedFollow(ctx, createFeedParams)
	if err != nil {
		return fmt.Errorf("error creating the feed follow: %w", err)
	}

	fmt.Printf("user '%s' is following %s", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func Following(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: following")
	}

	ctx := context.Background()

	// get user id
	currentUser := s.Cfg.CurrentUser
	allFeeds, err := s.Db.GetFeedFollowsForUser(ctx, currentUser)
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}

	if len(allFeeds) == 0 {
		fmt.Println("The current user is not following any feeds.")
		return nil
	}

	fmt.Println("Following feeds:")
	for _, feed := range allFeeds {
		fmt.Printf("   - %s\n", feed.FeedName)
	}

	return nil
}

func Unfollow(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: unfollow <url>")
	}

	feedUrl := cmd.Args[0]

	ctx := context.Background()

	// check for url
	_, err := s.Db.GetFeed(ctx, feedUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("feed '%s' not found", feedUrl)
		}
		return err
	}

	// get user id
	currentUser := s.Cfg.CurrentUser
	user, err := s.Db.GetUserByName(ctx, currentUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user '%s' not found", currentUser)
		}
		return err
	}

	unfollowParams := database.UnfollowFeedParams{
		UserID: user.ID,
		Url:    feedUrl,
	}

	err = s.Db.UnfollowFeed(ctx, unfollowParams)
	if err != nil {
		return fmt.Errorf("error unfollowing feed '%s': %w", feedUrl, err)
	}

	fmt.Printf("user %s unfollowed %s", user.Name, feedUrl)
	return nil
}
