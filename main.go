package main

import (
	//"fmt"
	"log"
	"os"

	"github.com/vinwong7/blogaggregator/internal/config"
)

type state struct {
	cfg_ptr *config.Config
}

func main() {
	configFile, _ := config.Read()

	new_state := state{
		cfg_ptr: &configFile,
	}

	new_commands := commands{
		commands_map: make(map[string]func(*state, command) error),
	}

	new_commands.register("login", handlerLogin)

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
