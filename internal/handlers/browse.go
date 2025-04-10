package handlers

import (
	"fmt"
	"strconv"

	"github.com/lealre/gator/internal/commands"
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

	fmt.Print(limit)

	// get user ID

	// Seacrh posts

	// log posts

	return nil
}
