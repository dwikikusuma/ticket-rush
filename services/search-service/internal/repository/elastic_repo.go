package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/dwikikusuma/ticket-rush/services/search-service/internal/domain"
	"github.com/elastic/go-elasticsearch/v8"
)

type elasticRepo struct {
	client *elasticsearch.Client
}

// NewElasticRepo creates a new instance
func NewElasticRepo(client *elasticsearch.Client) domain.SearchRepository {
	return &elasticRepo{client: client}
}

func (r *elasticRepo) SearchQuery(query string, limit int, cursor string) (*domain.SearchResult, error) {
	var shouldQuery []map[string]interface{}

	if query != "" {
		shouldQuery = append(shouldQuery, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     query,
				"fields":    []string{"event_name", "stadium"},
				"fuzziness": "AUTO",
			},
		})
	} else {
		shouldQuery = append(shouldQuery, map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	}

	esQuery := map[string]interface{}{
		"size": limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": shouldQuery,
				"filter": []map[string]interface{}{
					{"term": map[string]interface{}{"status.keyword": "AVAILABLE"}},
				},
			},
		},
		"sort": []map[string]interface{}{
			{"id": "asc"},
		},
	}

	if cursor != "" {
		cursorInt, _ := strconv.Atoi(cursor)
		esQuery["search_after"] = []interface{}{cursorInt}
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(esQuery); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex("tickets"),
		r.client.Search.WithBody(&buf),
	)
	log.Println(res)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	log.Println(buf)
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	var tickets []domain.Ticket
	var nextCursor string

	for _, hit := range hits {
		h := hit.(map[string]interface{})
		source := h["_source"].(map[string]interface{})

		id := int(source["id"].(float64))

		tickets = append(tickets, domain.Ticket{
			ID:        id,
			EventName: source["event_name"].(string),
			Stadium:   source["stadium"].(string),
			Price:     int(source["price"].(float64)),
			SeatID:    source["seat_id"].(string),
			Status:    source["status"].(string),
		})

		if sortArr, ok := h["sort"].([]interface{}); ok && len(sortArr) > 0 {
			val := sortArr[0].(float64)
			nextCursor = fmt.Sprintf("%.0f", val)
		}
	}

	return &domain.SearchResult{
		Tickets:    tickets,
		NextCursor: nextCursor,
	}, nil
}
