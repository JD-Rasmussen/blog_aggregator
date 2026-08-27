package main

import (
	"fmt"

	"github.com/jd-rasmussen/blog_aggregator/internal/config"
)

type state struct {
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
	s.cfg.SetUser(username)
	fmt.Println("Current user:", username)
	return nil

}
