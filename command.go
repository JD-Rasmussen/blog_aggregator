package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jd-rasmussen/blog_aggregator/internal/config"
	"github.com/jd-rasmussen/blog_aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	Name string
	args []string
}

type commands struct {
	cmdHandler map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.cmdHandler[cmd.Name]
	if !ok {
		fmt.Println("unknown command")
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdHandler[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No username provided")
	}
	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("Error checking if user exists: %v", err)
	}
	s.cfg.SetUser(username)
	fmt.Println("Current user:", username)
	return nil
}

func handlerRegister(s *state, cmd command) error { // add new user to the database
	if len(cmd.args) < 1 {
		return fmt.Errorf("No username provided")
	}
	username := cmd.args[0]
	//check if user already exists in the database
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("Error checking if user exists: %v", err)
	}
	// Add user to the database
	_, err = s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("Error adding user to database: %v", err)
	}
	fmt.Println("User registered:", username)

	handlerLogin(s, command{Name: "login", args: []string{username}}) // Log in the user after registration

	return nil
}

func handlerReset(s *state, cmd command) error {
	// Reset the database
	err := s.db.ResetDatabase(context.Background())
	if err != nil {
		return fmt.Errorf("Error resetting database: %v", err)
	}
	fmt.Println("Database reset successfully")
	return nil
}

func handlerGetUsers(s *state, cmd command) error { // return all users in the database
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error fetching users: %v", err)
	}
	for _, user := range users {
		if s.cfg.Current_user_name == user.Name {
			fmt.Println(user.Name, "(current)")
			continue
		}
		fmt.Println("User:", user.Name)
	}
	return nil
}

func HandleAggregateFeeds(s *state, cmd command) error {
	//fetch rss feeds from url
	feedURL := "https://www.wagslane.dev/index.xml" //cmd.args[0]

	feed, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Error fetching feed: %v", err)
	}

	fmt.Println("Feed Title:", feed.Channel.Title)
	fmt.Println("Feed Link:", feed.Channel.Link)
	fmt.Println("Feed Description:", feed.Channel.Description)

	for _, item := range feed.Channel.Item {
		fmt.Println("Item Title:", item.Title)
		fmt.Println("Item Link:", item.Link)
		fmt.Println("Item Description:", item.Description)
		fmt.Println("Item PubDate:", item.PubDate)
	}

	return nil

}

func HandleAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No feed URL provided")
	}
	feedname := cmd.args[0]
	feedurl := cmd.args[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("Error fetching current user: %v", err)
	}

	// Add feed to the database
	_, err = s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedname,
		Url:       feedurl,
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("Error adding feed to database: %v", err)
	}
	fmt.Println("Feed added:", feedname, "with URL:", feedurl)

	HandleFollowFeed(s, command{Name: "follow", args: []string{feedurl}}) // Follow the feed after adding it

	return nil
}

func HandleFeeds(s *state, cmd command) error {
	// Fetch all feeds in the database
	feeds, err := s.db.GetFeedsWithUser(context.Background())
	if err != nil {
		return fmt.Errorf("Error fetching feeds: %v", err)
	}

	for _, feed := range feeds {
		fmt.Println("Feed Name:", feed.Name)
		fmt.Println("Feed URL:", feed.Url)
		fmt.Println("Feed Created by:", feed.UserName)
	}

	return nil
}

func HandleFollowFeed(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No feed name provided")
	}
	feedurl := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedurl)
	if err != nil {
		return fmt.Errorf("Error fetching feed: %v", err)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("Error fetching current user: %v", err)
	}

	// Add feed follow to the database
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		FeedID:    uuid.NullUUID{UUID: feed.ID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("Error following feed: %v", err)
	}
	fmt.Println("User", user.Name, "is now following feed:", feed.Name, "with URL:", feed.Url)

	return nil
}

func HandleFollowing(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	if err != nil {
		return fmt.Errorf("Error fetching current user: %v", err)
	}

	// Fetch all feeds followed by the current user
	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), uuid.NullUUID{UUID: user.ID, Valid: true})
	if err != nil {
		return fmt.Errorf("Error fetching followed feeds: %v", err)
	}

	for _, feedFollow := range followedFeeds {
		fmt.Println("Feed Name:", feedFollow.FeedName)
		fmt.Println("Feed URL:", feedFollow.FeedID) // Assuming you want to print the Feed ID here
		fmt.Println("Followed by User:", feedFollow.UserName)
	}

	return nil
}
