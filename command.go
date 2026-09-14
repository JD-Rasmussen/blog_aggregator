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
