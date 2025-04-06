package handlers

import (
	"context"
	"fmt"

	"github.com/lealre/gator/internal/commands"
)

func ListUsers(s *commands.State, cmd commands.Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("usage: reset")
	}
	ctx := context.Background()

	usersList, err := s.Db.ListUsers(ctx)
	if err != nil {
		return fmt.Errorf("error reseting database: %w", err)
	}

	if len(usersList) == 0 {
		fmt.Print("No users are registered")
		return nil
	}

	for _, user := range usersList {
		if s.Cfg.CurrentUser == user.Name {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Println(user.Name)
		}
	}

	return nil
}
