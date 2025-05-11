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

func handlerReset(s *state, cmd command) error {

	err := s.db.Reset(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("User database has been reset.")
	return nil
}

func handlerUserList(s *state, cmd command) error {

	userList, err := s.db.UserList(context.Background())
	if err != nil {
		return err
	}

	for _, user := range userList {
		if user == s.cfg_ptr.Current_user_name {
			fmt.Printf("%v (current)\n", user)
		} else {
			fmt.Printf("%v\n", user)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {

	feedStruct, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", *feedStruct)

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {

	if len(cmd.arguments) < 2 {
		log.Fatal("Missing name of feed or feed URL. Exiting...\n")

	}
	/*
		userInfo, err := s.db.GetUser(context.Background(), s.cfg_ptr.Current_user_name)
		if err != nil {
			log.Fatal(err)
		}
	*/
	s.db.CreateFeed(context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.arguments[0],
			Url:       cmd.arguments[1],
			UserID:    user.ID,
		},
	)

	feedData, err := s.db.GetFeed(context.Background(), cmd.arguments[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Feed has been added to database. Here is the entry.")
	fmt.Printf("Name: %v, URL: %v\n", feedData.Name, feedData.Url)
	fmt.Printf("Created_At: %v, UserID: %v\n", feedData.CreatedAt, feedData.UserID)

	s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feedData.ID,
		})

	return nil
}

func handlerfeedList(s *state, cmd command) error {

	feedList, err := s.db.FeedList(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feedList {
		fmt.Printf("Feed Name: %v, Feed URL: %v, Added By: %v\n", feed.Feedname, feed.Url, feed.Username)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {

	if len(cmd.arguments) < 1 {
		log.Fatal("Missing feed URL. Exiting...\n")

	}

	/*
		userInfo, err := s.db.GetUser(context.Background(), s.cfg_ptr.Current_user_name)
		if err != nil {
			log.Fatal(err)
		}
	*/

	feedInfo, err := s.db.GetFeed(context.Background(), cmd.arguments[0])
	if err != nil {
		log.Fatal(err)
	}

	feedFollowData, err := s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feedInfo.ID,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current user, %v, now following feed, %v\n", feedFollowData.UserName, feedFollowData.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	feedFollowData, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		log.Fatal(err)
	}

	println("You are following these feeds currently:")
	for _, feed := range feedFollowData {
		fmt.Printf("%v\n", feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		log.Fatal("Missing feed URL. Exiting...\n")

	}

	feedInfo, err := s.db.GetFeed(context.Background(), cmd.arguments[0])
	if err != nil {
		log.Fatal(err)

	}

	err = s.db.Unfollow(context.Background(),
		database.UnfollowParams{
			FeedID: feedInfo.ID,
			UserID: user.ID,
		})
	if err != nil {
		log.Fatal(err)

	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {
		userInfo, err := s.db.GetUser(context.Background(), s.cfg_ptr.Current_user_name)
		if err != nil {
			log.Fatal(err)
		}
		return handler(s, cmd, userInfo)

	}
}
