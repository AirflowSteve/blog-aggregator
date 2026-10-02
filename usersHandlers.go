package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/AirflowSteve/blog-aggregator/internal/database"
)

func handlerLogin(s *State, cmd command) error {
	if len(cmd.Args) != 1 {
		fmt.Println("login handler expects a username, none found")
		os.Exit(1)
	}

	usr_name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), usr_name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.config.SetUser(usr_name)
	if err != nil {
		return err
	}

	fmt.Println("User logged in successfully!")
	return nil
}

func handlerRegister(s *State, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)

	}

	usr_name := cmd.Args[0]

	usr, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      usr_name,
	})
	if err != nil {
		return err
	}

	err = s.config.SetUser(usr.Name)
	if err != nil {
		return err
	}
	printUser(usr)
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:		%v\n", user.ID)
	fmt.Printf(" * Name:	%v\n", user.Name)
}
func handlerReset(s *State, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("could not reset the table: %w\n", err)
	}
	s.config.SetUser("")
	fmt.Println("the table was reset successfully!")
	return nil
}

func handlerGetUsers(s *State, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't select the users: %w\n", err)
	}
	for _, usr := range users {
		if usr.Name == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", usr.Name)
		} else {
			fmt.Printf("* %s\n", usr.Name)
		}
	}
	return nil
}
