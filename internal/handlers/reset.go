package handlers

import (
	"context"
	"fmt"

	"github.com/lealre/gator/internal/commands"
)

func Reset(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: reset")
	}
	ctx := context.Background()

	err := s.Db.ResetTable(ctx)
	if err != nil {
		return fmt.Errorf("error reseting database: %w", err)
	}

	fmt.Println("Users table reset")
	return nil
}
