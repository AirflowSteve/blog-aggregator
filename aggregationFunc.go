package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/AirflowSteve/blog-aggregator/internal/database"
	"github.com/google/uuid"
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
		t, err := time.Parse(time.RFC1123Z, item.PubDate)
		published_time := sql.NullTime{}
		if err == nil {
			published_time.Time = t
			published_time.Valid = true

		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: sql.NullString{
				String: item.Title,
				Valid:  true,
			},
			Url: sql.NullString{
				String: item.Link,
				Valid:  true,
			},
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: published_time,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			fmt.Printf("Couldn't create post: %v", err)
			continue
		}

	}
	fmt.Printf("Feed %s collected, %v posts found\n", feed.Name, len(items))
	return nil
}
