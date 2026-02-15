package events

import "time"

const (
	BookingEvents = "booking.created"
)

type BookingCreatedEvent struct {
	BookingID string    `json:"booking_id"`
	UserID    string    `json:"user_id"`
	EventName string    `json:"event_name"`
	SeatID    string    `json:"seat_id"`
	Amount    int64     `json:"amount"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}

type DlqMessage struct {
	Payload     []byte    `json:"payload"`
	ErrorReason string    `json:"error_reason"`
	ErrorDetail error     `json:"error_detail"`
	FailedAt    time.Time `json:"failed_at"`
	Service     string    `json:"service"`
}
