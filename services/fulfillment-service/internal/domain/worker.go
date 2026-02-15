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
	consumer    *events.Consumer
	db          ReceiptRepo
	dlqProducer *events.Producer
}

func NewConsumerWorker(consumer *events.Consumer, db ReceiptRepo, dlqProducer *events.Producer) *ConsumerWorker {
	return &ConsumerWorker{
		consumer:    consumer,
		db:          db,
		dlqProducer: dlqProducer,
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
			w.handleFailedMessage(ctx, mess, "JSON_UNMARSHAL_ERROR", decodeErr)
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
			w.handleFailedMessage(ctx, mess, "DB_SAVE_ERROR", err)
			continue
		}

		if err = w.consumer.CommitMessages(ctx, mess); err != nil {
			log.Printf("Error committing message: %v", err)
			w.handleFailedMessage(ctx, mess, "KAFKA_COMMIT_ERROR", err)
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

func (w *ConsumerWorker) handleFailedMessage(ctx context.Context, msg kafka.Message, reason string, err error) {
	log.Println("Handling failed message will be implemented later")
	dlqMsg := events.DlqMessage{
		Payload:     msg.Value,
		ErrorReason: reason,
		ErrorDetail: err,
		FailedAt:    time.Now(),
		Service:     "fulfillment-service",
	}

	dlqData, marshalErr := json.Marshal(dlqMsg)
	if marshalErr != nil {
		log.Printf("Error marshaling DLQ message: %v", marshalErr)
	}

	if pubErr := w.dlqProducer.Publish(ctx, "dlq.booking", "", dlqData); pubErr != nil {
		log.Printf("Error publishing DLQ message: %v", pubErr)
	}

	if err = w.consumer.CommitMessages(ctx, msg); err != nil {
		log.Printf("Error committing failed message: %v", err)
	}
}
