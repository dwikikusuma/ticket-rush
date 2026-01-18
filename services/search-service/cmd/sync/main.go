package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

const (
	postgresURL = "postgres://user:password@localhost:5432/ticket_db"
	elasticURL  = "http://localhost:9200"
	indexName   = "tickets"
)

type Ticket struct {
	ID        int       `json:"id"`
	EventName string    `json:"event_name"`
	Stadium   string    `json:"stadium"`
	Price     int       `json:"price"`
	SeatID    string    `json:"seat_id"`
	Status    string    `json:"status"`
	EventDate time.Time `json:"event_date"`
}

func main() {
	log.Println(" Starting synchronization service...")
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, postgresURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{elasticURL},
	})
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %v", err)
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  indexName,
		Client: es,
	})
	if err != nil {
		log.Fatalf("Error creating the bulk indexer: %v", err)
	}

	rows, err := conn.Query(ctx, "SELECT id, event_name, stadium, price, seat_id, status, event_date FROM tickets;")
	if err != nil {
		log.Fatalf("Error querying tickets: %v", err)
	}
	defer rows.Close()

	start := time.Now()
	total := 0

	for rows.Next() {
		var t Ticket
		err = rows.Scan(&t.ID, &t.EventName, &t.Stadium, &t.Price, &t.SeatID, &t.Status, &t.EventDate)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		data, err := json.Marshal(t)
		if err != nil {
			log.Fatalf("Error marshaling ticket to JSON: %v", err)
			continue
		}

		err = bi.Add(ctx, esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: fmt.Sprintf("%d", t.ID),
			Body:       bytes.NewReader(data),
		})

		if err != nil {
			log.Fatalf("Unexpected error while adding item to bulk indexer: %v", err)
		}

		total++
		if total%5000 == 0 {
			log.Printf("  Indexed %d tickets...", total)
		}
	}
	if err = bi.Close(ctx); err != nil {
		log.Fatalf("Error closing the bulk indexer: %v", err)
	}

	biStats := bi.Stats()
	duration := time.Since(start)
	log.Printf("✅ Synchronization complete: %d tickets indexed in %s", biStats.NumFlushed, duration)
}
