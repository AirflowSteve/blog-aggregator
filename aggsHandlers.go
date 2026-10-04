package main

import (
	"context"
	"fmt"
	"time"

	"github.com/AirflowSteve/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *State, cmd command) error {
	result, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", result)
	return nil
}

func handlerAddFeed(s *State, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("* usage: %v <name> <url>\n", cmd.Name)
	}

	ctxt := context.Background()

	usr, err := s.db.GetUser(ctxt, s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't retrieve the user's ID: %w", err)
	}

	feedName := cmd.Args[0]
	url := cmd.Args[1]

	// result, err := fetchFeed(ctxt, url)
	// if err != nil {
	// 	return fmt.Errorf("couldn't fetch the feed: %w\n", err)
	// }

	feed, err := s.db.CreateFeed(ctxt, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       url,
		UserID:    usr.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't insert new feed into the table", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("=====================================")

	return nil

}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}
