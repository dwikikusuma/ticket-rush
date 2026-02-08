package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	bookingURL = "http://localhost:8085/bookings"
	authURL    = "http://localhost:8087/login"

	// 🎯 TARGET: Everyone fights for this seat!
	eventName = "Concert: laudantium #1"
	seatID    = "Seat-11"

	totalUsers = 50
)

func main() {
	// 1. Setup: Ensure the ticket is AVAILABLE before we start
	fmt.Println("🔧 preparing test...")

	// 2. Get a valid JWT (Simulating a logged-in user)
	token, err := login()
	if err != nil {
		log.Fatalf("❌ Failed to login: %v", err)
	}
	log.Println("🔑 Got JWT Token")

	// 3. The Rush! 🏃‍♂️💨
	var wg sync.WaitGroup
	wg.Add(totalUsers)

	successCount := 0
	failCount := 0
	var mu sync.Mutex

	start := time.Now()
	log.Printf("🚀 Starting %d concurrent booking requests for Seat %s...", totalUsers, seatID)

	for i := 0; i < totalUsers; i++ {
		go func(id int) {
			defer wg.Done()

			// Attempt to book
			status := bookTicket(token, eventName, seatID)

			mu.Lock()
			if status == 200 {
				successCount++
				fmt.Printf("User %d: ✅ WON THE TICKET!\n", id)
			} else {
				failCount++
				fmt.Printf("User %d: ❌ Failed (%d)\n", id, status)
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	// 4. The Verdict
	fmt.Println("\n--- 🏁 STRESS TEST RESULTS ---")
	fmt.Printf("Time Taken: %v\n", duration)
	fmt.Printf("Total Requests: %d\n", totalUsers)
	fmt.Printf("Successful Bookings: %d\n", successCount)
	fmt.Printf("Failed Bookings:     %d\n", failCount)

	if successCount == 1 {
		fmt.Println("\n✅ PASS: System correctly sold the ticket exactly once.")
	} else if successCount == 0 {
		fmt.Println("\n⚠️  WARNING: No one got the ticket? (Check logs)")
	} else {
		fmt.Printf("\n🚨 FAIL: CRITICAL! Double booking detected! Sold %d times.\n", successCount)
	}
}

// --- Helpers ---

func login() (string, error) {
	payload := map[string]string{
		"email":    "test@test.com",
		"password": "agustinnus",
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(authURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed: %s", string(bodyBytes))
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	return res["token"].(string), nil
}

func bookTicket(token, event, seat string) int {
	payload := map[string]string{
		"event_name": event,
		"seat":       seat,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", bookingURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
