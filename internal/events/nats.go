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

// Encode a specific message, to send information related to that
func (natsStore *NatsEventStore) encodeMessage(message Message) ([]byte, error) {
	b := bytes.Buffer{}

	err := gob.NewEncoder(&b).Encode(message) // Encode message into bytes
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil // Return message encoded
}

// Deocode bytes of data into a specific interfae type (kind must be Message interface)
func (natsStore *NatsEventStore) decodeMessage(data []byte, message interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	// This line is going to try to return the bytes into the message interface
	return gob.NewDecoder(&b).Decode(message)
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

func (natsStore *NatsEventStore) OnCreatedFeed(function func(CreatedFeedMessage)) (err error) {
	message := CreatedFeedMessage{}

	// We are going to trying to suscribe into a specific message
	natsStore.feedCreatedSub, err = natsStore.conn.Subscribe(message.Type(), func(msg *nats.Msg) {
		natsStore.decodeMessage(msg.Data, &message) // Decode message into Message type
		function(message)                           // Return the message to sended function
	})
	return
}

func (natsStore *NatsEventStore) SuscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	var err error
	message := CreatedFeedMessage{}

	natsStore.feedCreatedChan = make(chan CreatedFeedMessage, 64) // Channel for new feeds created
	ch := make(chan *nats.Msg, 64)                                // Channel for information from nats service (byte data)

	// Suscribre Nat service into a new channel
	natsStore.feedCreatedSub, err = natsStore.conn.ChanSubscribe(message.Type(), ch)
	if err != nil {
		return nil, err
	}

	// New concurrent methods to handle channel interactions
	go func() {
		for {
			select {
			case msg := <-ch: // When channel recives data
				natsStore.decodeMessage(msg.Data, &message) // Decode message
				natsStore.feedCreatedChan <- message        // Send message decoded into created feed channel
			}
		}
	}()

	return (<-chan CreatedFeedMessage)(natsStore.feedCreatedChan), nil
}
