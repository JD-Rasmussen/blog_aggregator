package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jd-rasmussen/blog_aggregator/internal/config"
	"github.com/jd-rasmussen/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

//postgresql://USERNAME:PASSWORD@HOST:PORT/DATABASE
//postgresql://postgres:1234@localhost:5432/gator
// psql -h localhost -U postgres gator
//goose -dir sql/schema postgres postgresql://postgres:1234@localhost:5432/gator up
//goose -dir sql/schema postgres postgresql://postgres:1234@localhost:5432/gator down

func main() {
	cmds := commands{cmdHandler: make(map[string]func(*state, command) error)} // Initialize the commands struct with an empty map

	cfg, err := config.Read() // Read the configuration from the file
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.Db_url) // Open a connection to the PostgreSQL database
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	states := state{db: dbQueries, cfg: &cfg} // Initialize state with the config

	cmds.register("login", handlerLogin)       // Register the login command
	cmds.register("register", handlerRegister) // Register the register command
	cmds.register("reset", handlerReset)       // Register the reset command

	if len(os.Args) < 2 { //check if a command is presented
		fmt.Println("No command provided")
		os.Exit(1)
	}

	err = cmds.run(&states, command{Name: os.Args[1], args: os.Args[2:]}) // Run the command based on command-line arguments})
	if err != nil {
		fmt.Println("Error running command:", err)
		os.Exit(1)
	}

	fmt.Println("current user:", states.cfg.Current_user_name)
}
