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
const totalRows = 1_000_000 // 1 Million

func main() {
	ctx := context.Background()

	// 1. Connect
	fmt.Println("🔌 Connecting to Database...")
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	// 2. Clear existing data
	fmt.Println("🧹 Clearing old data...")
	_, err = conn.Exec(ctx, "TRUNCATE TABLE tickets;")
	if err != nil {
		log.Printf("Warning: Could not truncate table: %v", err)
	}

	// 3. Bulk Insert
	fmt.Printf("🚀 Starting Seed of %d rows...\n", totalRows)
	startTime := time.Now()

	rowsGenerated := 0

	// --- STATE VARIABLES ---
	var (
		currentEventName    string
		currentStadium      string
		currentEventDate    time.Time
		currentPrice        int
		seatsInCurrentEvent int
		currentSeatIndex    int
		eventGlobalCounter  int // NEW: Tracks how many events we've created total
	)

	// Initialize to force creation on first loop
	seatsInCurrentEvent = 0
	currentSeatIndex = 0
	eventGlobalCounter = 0

	count, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"tickets"},
		[]string{"event_name", "stadium", "price", "seat_id", "status", "event_date"},
		pgx.CopyFromFunc(func() ([]any, error) {
			if rowsGenerated >= totalRows {
				return nil, nil
			}

			// --- CHECK IF WE NEED A NEW EVENT ---
			if currentSeatIndex >= seatsInCurrentEvent {
				eventGlobalCounter++ // Increment unique event ID

				// 1. Generate new Event details with UNIQUE ID
				// "Concert: Rock #1", "Concert: Jazz #2", etc.
				// This prevents collision even if Faker returns the same word twice.
				currentEventName = fmt.Sprintf("Concert: %s #%d", faker.Word(), eventGlobalCounter)

				currentStadium = "Stadium " + faker.Word()
				currentEventDate = time.Now().AddDate(0, 0, rand.Intn(90))

				// 2. Randomize seats for this event (between 100 and 2000)
				seatsInCurrentEvent = rand.Intn(1900) + 100

				// 3. Set Price
				currentPrice = rand.Intn(200000) + 50000

				// 4. Reset seat counter
				currentSeatIndex = 0
			}

			currentSeatIndex++
			rowsGenerated++

			if rowsGenerated%100000 == 0 {
				fmt.Printf("   ... %d rows generated\n", rowsGenerated)
			}

			return []any{
				currentEventName,
				currentStadium,
				currentPrice,
				fmt.Sprintf("Seat-%d", currentSeatIndex),
				"AVAILABLE",
				currentEventDate,
			}, nil
		}),
	)

	if err != nil {
		log.Fatalf("❌ Seeding failed: %v", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("\n✅ SUCCESS! Inserted %d rows in %.2f seconds.\n", count, duration.Seconds())
}
