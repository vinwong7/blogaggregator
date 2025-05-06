package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/vinwong7/blogaggregator/internal/config"
	"github.com/vinwong7/blogaggregator/internal/database"
)

type state struct {
	db      *database.Queries
	cfg_ptr *config.Config
}

func main() {

	configFile, err := config.Read()
	if err != nil {
		log.Fatal("Error reading config file. Exiting...")
	}

	db, err := sql.Open("postgres", configFile.Db_url)
	if err != nil {
		log.Fatal("Error generating database. Exiting...")
	}

	dbQueries := database.New(db)

	new_state := state{
		db:      dbQueries,
		cfg_ptr: &configFile,
	}

	new_commands := commands{
		commands_map: make(map[string]func(*state, command) error),
	}

	new_commands.register("login", handlerLogin)
	new_commands.register("register", handlerRegister)
	new_commands.register("reset", handlerReset)
	new_commands.register("users", handlerUserList)

	if len(os.Args) < 2 {
		log.Fatal("Not enought arguments provided. Exiting...\n")
	}

	command_name := os.Args[1]
	command_args := os.Args[2:]

	command_instance := command{
		name:      command_name,
		arguments: command_args,
	}

	new_commands.run(&new_state, command_instance)

}
