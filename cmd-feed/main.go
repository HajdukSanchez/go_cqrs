package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hajduksanchez/go_cqrs/internal/database"
	"github.com/hajduksanchez/go_cqrs/internal/events"
	"github.com/hajduksanchez/go_cqrs/internal/handlers"
	"github.com/hajduksanchez/go_cqrs/internal/repository"
	"github.com/kelseyhightower/envconfig"
)

// Env keys defined in env file
type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Create postgres connection
	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresDB)
	// Create new repository
	repo, err := database.NewFeedDataBase(addr)
	if err != nil {
		log.Fatal(err)
	}
	repository.SetRepository(repo)

	// Create nats connection
	nats := fmt.Sprintf("nats://%s", config.NatsAddress)
	// Create new nats connection
	eventStore, err := events.NewNatsEventStore(nats)
	if err != nil {
		log.Fatal(err)
	}
	events.SetEventStore(eventStore)
	// Close nats connection at the end
	defer events.Close()

	// Create new router for server
	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

// New mux router
func newRouter() (router *mux.Router) {
	router = mux.NewRouter() // Create router

	// Add new handler for routes
	router.HandleFunc("/feeds", handlers.CreatedFeedHandler).Methods(http.MethodPost)
	return
}
