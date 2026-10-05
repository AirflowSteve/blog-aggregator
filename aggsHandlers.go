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

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    usr.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("couldn't create a follow: %w\n", err)
	}

	fmt.Println("Feed created successfully:")

	fmt.Println()
	fmt.Println("=====================================")
	printFeed(feed)
	fmt.Println()
	fmt.Println("Feed followed successfully:")

	fmt.Println()
	fmt.Println("=====================================")
	fmt.Printf("* Name of the feed: %s\n", follow.FeedName)
	fmt.Printf("* User: %s\n", follow.UserName)

	return nil

}

func handlerFeeds(s *State, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't retrieve feeds from the database: %w\n", err)
	}

	for _, feed := range feeds {
		printFeedRows(feed)
	}

	return nil
}

func handlerFollow(s *State, cmd command) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("* usage: %v <url>\n", cmd.Name)
	}
	url := cmd.Args[0]

	usr_id, err := s.db.GetUserByName(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get user's ID: %w\n", err)
	}

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't retrieve the feed from the database: %w\n", err)
	}

	feed_id := feed.ID

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    usr_id,
		FeedID:    feed_id,
	})

	if err != nil {
		return fmt.Errorf("couldn't create a follow: %w\n", err)
	}

	fmt.Println("Follow created successfully:")

	fmt.Println()
	fmt.Println("=====================================")
	fmt.Printf("* Name of the feed: %s\n", follow.FeedName)
	fmt.Printf("* User: %s\n", follow.UserName)

	return nil
}

func handlerFollowing(s *State, cmd command) error {
	usr_id, err := s.db.GetUserByName(context.Background(), s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't retrieve the current's user id: %w\n", err)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), usr_id)
	if err != nil {
		return fmt.Errorf("couldn't get the list of follows for the current user: %w\n", err)
	}

	fmt.Printf("* user: %s\n", s.config.CurrentUserName)
	fmt.Println()

	for _, follow := range follows {
		fmt.Printf("* Feed: %s\n", follow.FeedName)
	}

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

func printFeedRows(f database.GetFeedsRow) {
	fmt.Printf("* Name: 			%s\n", f.Name)
	fmt.Printf("* URL: 				%s\n", f.Url)
	fmt.Printf("* Created by:		%s\n", f.UsersName)
}
