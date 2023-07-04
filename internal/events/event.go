package events

import (
	"context"

	"github.com/hajduksanchez/go_cqrs/internal/models"
)

type EventStore interface {
	Close()
	PublishCreatedFeed(ctx context.Context, feed *models.Feed) error
	SuscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error)
	OnCreatedFeed(function func(CreatedFeedMessage)) error
}

var _eventStore EventStore

// Event to close connection with the store
func Close() {
	_eventStore.Close()
}

// Event to publish when a new feed is created
func PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	return _eventStore.PublishCreatedFeed(ctx, feed)
}

// Event to suscribe when a new feed is created
func SuscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	return _eventStore.SuscribeCreatedFeed(ctx)
}

// Event when a new feed is created
func OnCreatedFeed(function func(CreatedFeedMessage)) error {
	return _eventStore.OnCreatedFeed(function)
}