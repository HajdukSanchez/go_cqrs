package repository

import (
	"context"

	"github.com/hajduksanchez/go_cqrs/internal/models"
)

type Repository interface {
	Close()
	InsertFeed(ctx context.Context, feed *models.Feed) error
	ListFeeds(ctx context.Context) ([]*models.Feed, error)
}

var repository Repository

// Create a new instance of the repository to works as dependency injection
func SetRepository(newRepository Repository) {
	repository = newRepository
}

// Function to close repository connection
func Close() {
	repository.Close()
}

// Insertion of a new feeds into DB
func InsertFeed(ctx context.Context, feed *models.Feed) error {
	return repository.InsertFeed(ctx, feed)
}

// List all the feeds of the DB
func ListFeeds(ctx context.Context) ([]*models.Feed, error) {
	return repository.ListFeeds(ctx)
}
