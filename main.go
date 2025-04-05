package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/lealre/gator/internal/config"
	"github.com/lealre/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Print(err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Print(err)
	}

	dbQueries := database.New(db)
	s := &state{cfg: &cfg, db: dbQueries}

	cmd := &commands{commands: make(map[string]func(*state, command) error)}
	cmd.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Please provide a command")
		os.Exit(1)
	}

	if len(os.Args) < 3 {
		fmt.Println("Usage: login <name>")
		os.Exit(1)
	}
	userCmd := os.Args[1]
	args := os.Args[2:]

	command := command{name: userCmd, args: args}
	err = cmd.run(s, command)
	if err != nil {
		errorMessage := fmt.Errorf("error executing command %s: %w", userCmd, err)
		fmt.Println(errorMessage)
		os.Exit(1)
	}

}

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("missing username")
	}

	userName := cmd.args[0]
	err := s.cfg.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("User setted as %s\n", userName)
	return nil
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if f, ok := c.commands[cmd.name]; ok {
		f(s, command{name: cmd.name, args: cmd.args})
	} else {
		fmt.Printf("Command %s not found\n", cmd.name)
	}
	return nil
}
