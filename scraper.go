package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rubacodes/boot-dev-rss/internal/database"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scapring on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest) // atimer that emits into its channel after the specified duration
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency)) // context.Background is the global context, accessible everywhere
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}
		// simile al concetto di Task.when all , ma l'esecuzione delle routine e' concorrente (pari al numero di feed da cercare)
		// il wg.Wait attende il w.done di tutte le n operazioni concorrenti
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(wg, db, feed)
		}
		wg.Wait()
	}

}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeddAsFetch(context.Background(), feed.ID)
	if err != nil {
		log.Println("error marking feed as fetched:", err)
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("error fetching feed", err)
		return
	}
	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("couldn't parse date %v with err %v", item.PubDate, err)
			continue
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			PublishedAt: pubAt,
			Title:       item.Title,
			Description: description,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			// log only when the erro is not duplicate key
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("failed to create post:", err)
		}
	}
	log.Printf("Feed %s collected, %v post found", feed.Name, len(rssFeed.Channel.Item))
}
