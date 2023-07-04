package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/hajduksanchez/go_cqrs/internal/models"
)

type FeedDatabase struct {
	db *sql.DB
}

// Create a new instance of Feed Database
func NewFeedDataBase(url string) (*FeedDatabase, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &FeedDatabase{db}, nil
}

func (repository *FeedDatabase) Close() {
	repository.db.Close()
}

func (repository *FeedDatabase) InsertFeed(ctx context.Context, feed *models.Feed) error {
	_, err := repository.db.ExecContext(ctx, "INSERT INTO feeds (id, title, description) VALUES ($1, $2, $3)", feed.Id, feed.Title, feed.Description)
	return err
}

func (repository *FeedDatabase) ListFeeds(ctx context.Context) ([]*models.Feed, error) {
	rows, err := repository.db.QueryContext(ctx, "SELECT id, title, description, created_at FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeds := []*models.Feed{}

	for rows.Next() {
		feed := &models.Feed{}
		if err := rows.Scan(&feed.Id, &feed.Title, &feed.Description, &feed.CreatedAt); err != nil {
			return nil, err
		}
		feeds = append(feeds, feed) // Insert new feed
	}
	return feeds, nil
}
