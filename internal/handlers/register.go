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

func Register(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: register <name>")
	}

	userName := cmd.Args[0]

	ctx := context.Background()

	// Check if user already exists
	_, err := s.Db.GetUserByName(ctx, userName)
	if err != nil {
		if err != sql.ErrNoRows {
			return fmt.Errorf("error getting user: %w", err)
		}
	}

	// Create new user
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}

	_, err = s.Db.CreateUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	err = s.Cfg.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("User setted as %s\n", userName)
	fmt.Printf("User created as %s\n", userName)
	return nil
}
