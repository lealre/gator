package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
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
	cmd.register("register", handlerRegister)
	cmd.register("reset", handlerReset)

	if len(os.Args) < 2 {
		fmt.Println("Please provide a command")
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
		return fmt.Errorf("usage: login <name>")
	}

	ctx := context.Background()

	userName := cmd.args[0]

	// Check if user already exists
	_, err := s.db.GetUserByName(ctx, userName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not foud: %w", err)
		}
	}

	err = s.cfg.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("User setted as %s\n", userName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: register <name>")
	}

	userName := cmd.args[0]

	ctx := context.Background()

	// Check if user already exists
	_, err := s.db.GetUserByName(ctx, userName)
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
		Name:      cmd.args[0],
	}

	_, err = s.db.CreateUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	err = s.cfg.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("User setted as %s\n", userName)
	fmt.Printf("User created as %s\n", userName)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: reset")
	}
	ctx := context.Background()

	err := s.db.ResetTable(ctx)
	if err != nil {
		return fmt.Errorf("error reseting database: %w", err)
	}

	fmt.Println("Users table reset")
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
		return f(s, command{name: cmd.name, args: cmd.args})
	} else {
		fmt.Printf("Command %s not found\n", cmd.name)
	}
	return nil
}
