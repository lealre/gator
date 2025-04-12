package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/lealre/gator/internal/commands"
	"github.com/lealre/gator/internal/config"
	"github.com/lealre/gator/internal/database"
	"github.com/lealre/gator/internal/handlers"

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
	s := &commands.State{Cfg: &cfg, Db: dbQueries}

	cmd := &commands.Commands{Commands: make(map[string]func(*commands.State, commands.Command) error)}
	cmd.Register("login", handlers.Login)
	cmd.Register("register", handlers.Register)
	cmd.Register("reset", handlers.Reset)
	cmd.Register("users", handlers.ListUsers)
	cmd.Register("agg", handlers.Aggregate)
	cmd.Register("addfeed", handlers.AddFeed)
	cmd.Register("feeds", handlers.ListFeed)
	cmd.Register("follow", handlers.Follow)
	cmd.Register("following", handlers.Following)
	cmd.Register("unfollow", handlers.Unfollow)
	cmd.Register("browse", handlers.Browse)

	if len(os.Args) < 2 {
		fmt.Println("Please provide a command")
		os.Exit(1)
	}

	userCmd := os.Args[1]
	args := os.Args[2:]

	command := commands.Command{Name: userCmd, Args: args}
	err = cmd.Run(s, command)
	if err != nil {
		errorMessage := fmt.Errorf("error executing command %s: %w", userCmd, err)
		fmt.Println(errorMessage)
		os.Exit(1)
	}

}
