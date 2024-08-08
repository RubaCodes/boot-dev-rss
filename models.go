package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/rubacodes/boot-dev-rss/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func DatabaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func DatabaseFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		Url:       dbFeed.Url,
		UserID:    dbFeed.UserID,
	}
}
func DatabaseFeedsToFeeds(dbFeed []database.Feed) []Feed {
	feeds := []Feed{}
	for _, item := range dbFeed {
		feeds = append(feeds, Feed{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			Name:      item.Name,
			Url:       item.Url,
			UserID:    item.UserID,
		})
	}
	return feeds
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FeedId    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
}

func DatabaseFeedsFollowstoFeedFollows(item database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        item.ID,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
		FeedId:    item.FeedID,
		UserID:    item.UserID,
	}
}

func DatabaseFeedsFollowsToFeedsFollow(dbFeed []database.FeedFollow) []FeedFollow {
	feeds := []FeedFollow{}
	for _, item := range dbFeed {
		feeds = append(feeds, DatabaseFeedsFollowstoFeedFollows(item))
	}
	return feeds
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description *string   `json:"description"` // il pointer se null ritona null sulla stringa
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func DatabasePostToPost(post database.Post) Post {
	var description *string
	if post.Description.Valid {
		description = &post.Description.String
	}
	return Post{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		PublishedAt: post.PublishedAt,
		Title:       post.Title,
		Description: description,
		Url:         post.Url,
		FeedID:      post.FeedID,
	}
}
func DatabasePoststoPosts(posts []database.Post) []Post {
	t := []Post{}
	for _, item := range posts {
		t = append(t, DatabasePostToPost(item))
	}
	return t
}
