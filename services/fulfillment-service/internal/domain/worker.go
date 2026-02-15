package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dwikikusuma/ticket-rush/common/pkg/events"
	"github.com/segmentio/kafka-go"
)

type ConsumerWorker struct {
	consumer *events.Consumer
	db       ReceiptRepo
}

func NewConsumerWorker(consumer *events.Consumer, db ReceiptRepo) *ConsumerWorker {
	return &ConsumerWorker{
		consumer: consumer,
		db:       db,
	}
}

func (w *ConsumerWorker) Start(ctx context.Context) {
	log.Println("Starting Consumer Worker")
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Consumer Worker")
			return
		default:
		}

		var eventData events.BookingCreatedEvent

		mess, err := w.consumer.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Context cancelled, stopping worker")
				return
			}
			log.Printf("Error fetching message: %v", err)
			continue
		}

		if decodeErr := json.Unmarshal(mess.Value, &eventData); decodeErr != nil {
			log.Printf("Error unmarshaling message: %v", decodeErr)
			if commErr := w.consumer.CommitMessages(ctx, mess); commErr != nil {
				log.Printf("Error committing message: %v", commErr)
			}
			continue
		}

		params := Receipt{
			BookingID: eventData.BookingID,
			UserID:    "",
			EventName: eventData.EventName,
			SeatID:    eventData.SeatID,
			QRCode:    fmt.Sprintf("code-%x", eventData.SeatID),
			Status:    "ISSUED",
		}

		if err = w.saveWithRetry(params); err != nil {
			log.Printf("Error saving receipt after retries: %v", err)
			w.handleFailedMessage(ctx, mess)
			continue
		}

		if err = w.consumer.CommitMessages(ctx, mess); err != nil {
			log.Printf("Error committing message: %v", err)
		} else {
			log.Printf("Successfully processed message for BookingID=%s", eventData.BookingID)
		}

	}
}

func (w *ConsumerWorker) saveWithRetry(data Receipt) error {
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err := w.db.Save(&data)
		if err != nil {
			log.Printf("Error saving receipt (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(1 * time.Second)
		} else {
			log.Println("Receipt saved successfully")
			return nil
		}
	}
	return errors.New("failed to save receipt after multiple attempts")
}

func (w *ConsumerWorker) handleFailedMessage(ctx context.Context, msg kafka.Message) {
	log.Println("Handling failed message will be implemented later")
	if err := w.consumer.CommitMessages(ctx, msg); err != nil {
		log.Printf("Error committing failed message: %v", err)
	}
}
