package handlers

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lealre/gator/internal/commands"
)

func Login(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: login <name>")
	}

	ctx := context.Background()

	userName := cmd.Args[0]

	// Check if user already exists
	_, err := s.Db.GetUserByName(ctx, userName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not foud: %w", err)
		}
	}

	err = s.Cfg.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("User setted as %s\n", userName)
	return nil
}
