package main

import (
	"errors"
	"fmt"
	"log"
)

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commands_map map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {

	if len(cmd.arguments) == 0 {
		log.Fatal("Username is required. Exiting...\n")

	}

	s.cfg_ptr.SetUser(cmd.arguments[0])

	fmt.Printf("User name has been set to %v.\n", cmd.arguments[0])

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	function, ok := c.commands_map[cmd.name]
	if !ok {
		return errors.New("command not found")
	}
	return function(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands_map[name] = f
}
