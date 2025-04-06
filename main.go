package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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
	cmd.register("users", handlerListUsers)
	cmd.register("agg", hendlerAggregate)

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

func handlerListUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: reset")
	}
	ctx := context.Background()

	usersList, err := s.db.ListUsers(ctx)
	if err != nil {
		return fmt.Errorf("error reseting database: %w", err)
	}

	if len(usersList) == 0 {
		fmt.Print("No users are registered")
		return nil
	}

	for _, user := range usersList {
		if s.cfg.CurrentUser == user.Name {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Println(user.Name)
		}
	}

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

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting %s: %w", feedUrl, err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feedData RSSFeed
	err = xml.Unmarshal(body, &feedData)
	if err != nil {
		return nil, err
	}

	decodeHTMLEntities(&feedData)

	return &feedData, nil
}

func decodeHTMLEntities(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
}

func hendlerAggregate(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()

	data, err := fetchFeed(ctx, url)
	if err != nil {
		return err
	}

	fmt.Print(data)
	return nil
}
