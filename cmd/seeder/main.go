package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5"
)

const dbURL = "postgres://user:password@localhost:5432/ticket_db"
const totalRows = 1_000_000

func main() {
	ctx := context.Background()

	// 1. Connect
	fmt.Println("🔌 Connecting to Database...")
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	// --- REMOVED: CREATE TABLE Logic ---
	// We assume 'make migrate-up' was run before this.

	// 2. Clear existing data (Optional: good for repeated tests)
	fmt.Println("🧹 Clearing old data...")
	_, err = conn.Exec(ctx, "TRUNCATE TABLE tickets;")
	if err != nil {
		log.Printf("Warning: Could not truncate table: %v", err)
	}

	// 3. Bulk Insert
	fmt.Printf("🚀 Starting Seed of %d rows...\n", totalRows)
	startTime := time.Now()
	rowsGenerated := 0

	count, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"tickets"},
		[]string{"event_name", "stadium", "price", "seat_id", "status", "event_date"},
		pgx.CopyFromFunc(func() ([]any, error) {
			if rowsGenerated >= totalRows {
				return nil, nil
			}
			rowsGenerated++
			if rowsGenerated%100000 == 0 {
				fmt.Printf("   ... %d rows\n", rowsGenerated)
			}
			return []any{
				"Concert: " + faker.Word(),
				"Stadium " + faker.Word(),
				rand.Intn(200000) + 50000,
				fmt.Sprintf("Seat-%d", rowsGenerated),
				"AVAILABLE",
				time.Now().AddDate(0, 0, rand.Intn(90)),
			}, nil
		}),
	)

	if err != nil {
		log.Fatalf("❌ Seeding failed: %v", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("\n✅ SUCCESS! Inserted %d rows in %.2f seconds.\n", count, duration.Seconds())
}
