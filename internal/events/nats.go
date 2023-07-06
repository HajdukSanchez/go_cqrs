package events

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/hajduksanchez/go_cqrs/internal/models"
	"github.com/nats-io/nats.go"
)

type NatsEventStore struct {
	conn            *nats.Conn              // Connection to nats
	feedCreatedSub  *nats.Subscription      // Suscribition to connect when a event is created
	feedCreatedChan chan CreatedFeedMessage // Connected channel
}

// Constructor to create a new Struct
func NewNatsEventStore(url string) (*NatsEventStore, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{
		conn: conn,
	}, nil
}

// Function to encode a specific message, to send information related to that
func (natsStore *NatsEventStore) encodeMessage(message Message) ([]byte, error) {
	b := bytes.Buffer{}

	err := gob.NewEncoder(&b).Encode(message) // Encode message into bytes
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil // Return message encoded
}

func (natsStore *NatsEventStore) Close() {
	if natsStore.conn != nil {
		natsStore.conn.Close() // Close nats connection
	}
	if natsStore.feedCreatedSub != nil {
		natsStore.feedCreatedSub.Unsubscribe() // Unsuscribre for feed
	}
	close(natsStore.feedCreatedChan) // Close channel connection
}

func (natsStore *NatsEventStore) PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	message := CreatedFeedMessage{
		Id:          feed.Id,
		Title:       feed.Title,
		Description: feed.Description,
		CreatedAt:   feed.CreatedAt,
	}

	// Encode new message to be send
	data, err := natsStore.encodeMessage(message)
	if err != nil {
		return err
	}

	// Publish a specific message with his data
	// This is how we are going to tell each connected microservice when there is a new feed
	return natsStore.conn.Publish(message.Type(), data)
}
