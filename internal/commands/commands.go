package commands

import (
	"fmt"

	"github.com/lealre/gator/internal/config"
	"github.com/lealre/gator/internal/database"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands map[string]func(*State, Command) error
}

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Commands[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	if f, ok := c.Commands[cmd.Name]; ok {
		return f(s, Command{Name: cmd.Name, Args: cmd.Args})
	} else {
		fmt.Printf("Command %s not found\n", cmd.Name)
	}
	return nil
}
