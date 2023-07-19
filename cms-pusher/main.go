package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hajduksanchez/go_cqrs/internal/events"
	"github.com/hajduksanchez/go_cqrs/internal/events/messages"
	"github.com/hajduksanchez/go_cqrs/internal/websocket"
	"github.com/kelseyhightower/envconfig"
)

// Env keys defined in env file
type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Creating new hub of clients
	hub := websocket.NewHub()

	// Create new Nats Store
	nats := fmt.Sprintf("nats://%s", config.NatsAddress)
	natsStore, err := events.NewNatsEventStore(nats)
	if err != nil {
		log.Fatal(err)
	}

	// We define a functioin to be handled when a new feed is created
	err = natsStore.OnCreatedFeed(func(m messages.CreatedFeedMessage) {
		// The hub of clients in the web socket, will be send a new message with the data of the new feed
		hub.Broadcast(NewCreatedFeedMessage(m.Id, m.Title, m.Description, m.CreatedAt), nil)
	})
	if err != nil {
		log.Fatal(err)
	}

	events.SetEventStore(natsStore)
	defer events.Close()

	// Start the web socket as a concurrency process
	go hub.Run()

	// We start the server with the websocket
	http.HandleFunc("/ws", hub.HandleWebSocket)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
