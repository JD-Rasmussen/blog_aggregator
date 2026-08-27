package main

import (
	"fmt"
	"os"

	"github.com/jd-rasmussen/blog_aggregator/internal/config"
)

func main() {
	cmds := commands{cmdHandler: make(map[string]func(*state, command) error)} // Initialize the commands struct with an empty map

	cfg, err := config.Read() // Read the configuration from the file
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}
	state := state{cfg: &cfg} // Initialize state with the config

	cmds.register("login", handlerLogin) // Register the login command

	if len(os.Args) < 2 { //check if a command is presented
		fmt.Println("No command provided")
		os.Exit(1)
	}

	err = cmds.run(&state, command{Name: os.Args[1], args: os.Args[2:]}) // Run the command based on command-line arguments})
	if err != nil {
		fmt.Println("Error running command:", err)
		os.Exit(1)
	}

}
