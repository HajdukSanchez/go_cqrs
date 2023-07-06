package search

import (
	"bytes"
	"context"
	"encoding/json"

	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/hajduksanchez/go_cqrs/internal/models"
)

type ElasticSearchRepository struct {
	client *elastic.Client
}

// Constructor
func NewElasticSearchRepository(url string) (*ElasticSearchRepository, error) {
	// Create a new elastic search instance
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}

	return &ElasticSearchRepository{client: client}, nil
}

func (repo *ElasticSearchRepository) Close() {
	// There is no way to close connection
}

func (repo *ElasticSearchRepository) IndexFeed(ctx context.Context, feed models.Feed) error {
	body, _ := json.Marshal(feed) // Feed json represetantion
	_, err := repo.client.Index(
		"feeds",               // Index Name
		bytes.NewReader(body), // Reader with data
		repo.client.Index.WithDocumentID(feed.Id), // Id of the document to e created
		repo.client.Index.WithContext(ctx),        // Context in case there is something wrong
		repo.client.Index.WithRefresh("wait_for"), // Refresh type defined in elastic search
	)
	return err
}
