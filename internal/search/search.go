package search

import (
	"context"

	"github.com/hajduksanchez/go_cqrs/internal/models"
)

type SearchRepository interface {
	Close()
	IndexFeed(ctx context.Context, feed models.Feed) error
	SearchFeed(ctx context.Context, query string) ([]models.Feed, error)
}

var _repository SearchRepository

// Set a new implementation for the repository
func SetSearchRepository(repository SearchRepository) {
	_repository = repository
}

// Close repository connection
func Close() {
	_repository.Close()
}

// Index a new feed
func IndexFeed(ctx context.Context, feed models.Feed) error {
	return _repository.IndexFeed(ctx, feed)
}

// Search feeds based on some query
func SearchFeed(ctx context.Context, query string) ([]models.Feed, error) {
	return _repository.SearchFeed(ctx, query)
}
