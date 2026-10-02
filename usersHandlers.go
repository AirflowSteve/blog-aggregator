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
		fmt.Errorf("usage: %v <name>", cmd.Name)

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
