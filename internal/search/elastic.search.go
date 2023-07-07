package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

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

func (repo *ElasticSearchRepository) SearchFeed(ctx context.Context, query string) (results []models.Feed, err error) {
	var buffer bytes.Buffer

	// Map with string keys and dynamic value
	// This is a json related to data we structure:
	// {
	// 	"query": {
	// 		"multi_match": {
	// 			"query": "abc",
	// 			"fields" [
	// 				"title",
	// 				"description"
	// 			],
	// 			"fuzziness": 3,
	// 			"cutoff_frequency": 0.0001
	// 		}
	// 	}
	// }
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":            query,
				"fields":           []string{"title", "description"},
				"fuzziness":        3,
				"cutoff_frequency": 0.0001,
			},
		},
	}
	// Encode search query into bytes
	if err = json.NewEncoder(&buffer).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := repo.client.Search(
		repo.client.Search.WithContext(ctx),         // Context for search
		repo.client.Search.WithIndex("feeds"),       // Index to search
		repo.client.Search.WithBody(&buffer),        // Data to search on indes
		repo.client.Search.WithTrackTotalHits(true), // Show how many matches found
	)
	if err != nil {
		return nil, err
	}

	// Validate at the end if there is an error closing response body
	defer func() {
		if err := res.Body.Close(); err != nil {
			results = nil
		}
	}()
	// If there is another error with data
	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var eRes map[string]interface{} // Variable to decode response and tranform into Feed model
	// Try to decode data and validate there is no error doing that
	if err = json.NewDecoder(res.Body).Decode(&eRes); err != nil {
		return nil, err
	}

	var feeds []models.Feed // Feeds to be send
	// This is similar to get key "hits" and tranform into a map, then inside this map,
	// we get a key "hits" and tranform into a slice
	// It look like this:
	// {
	// 	"hits": {
	// 		"hits": [...] // This is where we are going to range the for loop
	// 	}
	// }
	for _, hit := range eRes["hits"].(map[string]interface{})["hits"].([]interface{}) {
		feed := models.Feed{}
		// {
		// 	"hits": {
		// 		"hits": [
		// 			{
		// 				"_source": ...
		// 			},
		// 			...
		// 		]
		// 	}
		// }
		source := hit.(map[string]interface{})["_source"]
		marshal, err := json.Marshal(source) // Get bytes of the _source
		if err != nil {
			return nil, err
		}
		// Tranform data into data structure
		if err = json.Unmarshal(marshal, &feed); err == nil {
			feeds = append(feeds, feed)
		}
	}

	return feeds, nil
}
