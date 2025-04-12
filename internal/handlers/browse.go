package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/lealre/gator/internal/commands"
	"github.com/lealre/gator/internal/database"
)

const DefaultLimit = 2

func Browse(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: browse <limit> [optional, default is 2]")
	}

	limit := DefaultLimit

	if len(cmd.Args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			fmt.Printf("Error parsing limit: %v. Using default value: %d\n", err, DefaultLimit)
		} else {
			limit = parsedLimit
		}
	}

	ctx := context.Background()

	user, err := s.Db.GetUserByName(ctx, s.Cfg.CurrentUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not foud: %w", err)
		}
	}

	listPostParams := database.ListPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.Db.ListPostsByUser(ctx, listPostParams)
	if err != nil {
		return fmt.Errorf("error getting post from user: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title.String)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Printf("Description: %s\n", post.Description.String)
		fmt.Printf("FeedId: %s\n", post.FeedID)
	}

	return nil
}
