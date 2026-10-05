package main

import (
	"context"
	"fmt"
	"log"
)

func scrapeFeeds(s *State) error {

	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	log.Println("Found a feed to fetch!")

	_, err = s.db.MarkFeedFetched(context.Background(), feed.ID)

	if err != nil {
		return err
	}

	feeds, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	items := feeds.Channel.Item

	for _, item := range items {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	fmt.Printf("Feed %s collected, %v posts found\n", feed.Name, len(items))
	return nil
}
