package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vinwong7/blogaggregator/internal/database"
)

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commands_map map[string]func(*state, command) error
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

func handlerLogin(s *state, cmd command) error {

	if len(cmd.arguments) == 0 {
		log.Fatal("Username is required. Exiting...\n")

	}

	userCheck, _ := s.db.GetUser(context.Background(), cmd.arguments[0])
	if userCheck.Name == "" {
		log.Fatal("Username is not in database. Exiting...\n")
	}

	s.cfg_ptr.SetUser(cmd.arguments[0])

	fmt.Printf("User name has been set to %v.\n", cmd.arguments[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		log.Fatal("Username is required. Exiting...\n")

	}

	userCheck, _ := s.db.GetUser(context.Background(), cmd.arguments[0])
	if userCheck.Name != "" {
		log.Fatal("Username already exists. Exiting...\n")
	}

	s.db.CreateUser(context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.arguments[0]},
	)

	s.cfg_ptr.SetUser(cmd.arguments[0])

	fmt.Printf("User has been created. User name was %v.\n", cmd.arguments[0])

	return nil
}
