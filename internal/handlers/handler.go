package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/hajduksanchez/go_cqrs/internal/events"
	"github.com/hajduksanchez/go_cqrs/internal/models"
	"github.com/hajduksanchez/go_cqrs/internal/repository"
	"github.com/segmentio/ksuid"
)

// Expected request on create a new feed
type createdFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// / Handler for function on create a new feed
func CreatedFeedHandler(w http.ResponseWriter, r *http.Request) {
	var request createdFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		// Bad request if there is and error decoding data into struct
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAt := time.Now().UTC() // Set now time for value created at
	id, err := ksuid.NewRandom()  // Create new random UID
	if err != nil {
		// Internal error trying craeting Feed ID
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	feed := models.Feed{
		Id:          id.String(),
		Title:       request.Title,
		Description: request.Description,
		CreatedAt:   createdAt,
	}

	// Insert feed into DB
	if err := repository.InsertFeed(r.Context(), &feed); err != nil {
		// Internal error inserting new feed
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Publis new event
	if err := events.PublishCreatedFeed(r.Context(), &feed); err != nil {
		log.Printf("Failed to publish created feed event: %v", err)
	}

	// Send response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(feed)
}
