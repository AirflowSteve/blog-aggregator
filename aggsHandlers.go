package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/AirflowSteve/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *State, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("* usage %v <time between requests>\n", cmd.Name)
	}

	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", duration)

	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func handlerAddFeed(s *State, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("* usage: %v <name> <url>\n", cmd.Name)
	}

	ctxt := context.Background()

	feedName := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(ctxt, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't insert new feed into the table", err)
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
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

func handlerFollow(s *State, cmd command, user database.User) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("* usage: %v <url>\n", cmd.Name)
	}
	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't retrieve the feed from the database: %w\n", err)
	}

	feed_id := feed.ID

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
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

func handlerFollowing(s *State, cmd command, user database.User) error {

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
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

func handlerUnfollow(s *State, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("* usage %v <url>\n", cmd.Name)
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Feed unfollowed successfully!\n")
	return nil

}

func middlewareLoggedIn(handler func(s *State, cmd command, user database.User) error) func(s *State, cmd command) error {
	return func(s *State, cmd command) error {
		usr, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("couldn't retrieve the current's user id: %w\n", err)
		}

		return handler(s, cmd, usr)
	}
}

func handlerBrowse(s *State, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		l, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("* usage %v <limit int>\n", cmd.Name)
		} else {
			limit = l
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)

	for _, post := range posts {
		printPost(post)
	}

	return nil
}

func printPost(post database.GetPostsForUserRow) {
	fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
	fmt.Printf("--- %s ---\n", post.Title.String)
	fmt.Printf("    %v\n", post.Description.String)
	fmt.Printf("Link: %s\n", post.Url.String)
	fmt.Println("=====================================")

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
