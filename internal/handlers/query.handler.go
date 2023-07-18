package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/hajduksanchez/go_cqrs/internal/events/messages"
	"github.com/hajduksanchez/go_cqrs/internal/models"
	"github.com/hajduksanchez/go_cqrs/internal/repository"
	"github.com/hajduksanchez/go_cqrs/internal/search"
)

// Index new feed when something is created
func OnCreatedFeed(message messages.CreatedFeedMessage) {
	feed := models.Feed{
		Id:          message.Id,
		Description: message.Description,
		Title:       message.Title,
		CreatedAt:   message.CreatedAt,
	}

	if err := search.IndexFeed(context.Background(), feed); err != nil {
		log.Printf("Failed to index feed: %v", err)
	}
}

// Handler to lista all the feeds
func ListFeedHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := repository.ListFeeds(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}

// Handler to get a specific feeds that match with a specific query values
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q") // Get q parameter with his value from URL
	if len(query) == 0 {
		http.Error(w, "Query value (q) is required", http.StatusBadRequest)
		return
	}

	feeds, err := search.SearchFeed(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}
