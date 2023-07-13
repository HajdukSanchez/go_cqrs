package events

import (
	"context"

	"github.com/hajduksanchez/go_cqrs/internal/events/messages"
	"github.com/hajduksanchez/go_cqrs/internal/models"
)

type EventStore interface {
	Close()
	PublishCreatedFeed(ctx context.Context, feed *models.Feed) error
	SuscribeCreatedFeed(ctx context.Context) (<-chan messages.CreatedFeedMessage, error)
	OnCreatedFeed(function func(messages.CreatedFeedMessage)) error
}

var _eventStore EventStore

// Create a new instance of the evet store to works as dependency injection
func SetEventStore(eventStore EventStore) {
	_eventStore = eventStore
}

// Event to close connection with the store
func Close() {
	_eventStore.Close()
}

// Event to publish when a new feed is created
func PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	return _eventStore.PublishCreatedFeed(ctx, feed)
}

// Event to suscribe when a new feed is created
func SuscribeCreatedFeed(ctx context.Context) (<-chan messages.CreatedFeedMessage, error) {
	return _eventStore.SuscribeCreatedFeed(ctx)
}

// Event when a new feed is created
func OnCreatedFeed(function func(messages.CreatedFeedMessage)) error {
	return _eventStore.OnCreatedFeed(function)
}
