package domain

import (
	"time"
)

type Receipt struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	BookingID string    `bson:"booking_id" json:"booking_id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	EventName string    `bson:"event_name" json:"event_name"`
	SeatID    string    `bson:"seat_id" json:"seat_id"`
	QRCode    string    `bson:"qr_code" json:"qr_code"`
	IssuedAt  time.Time `bson:"issued_at" json:"issued_at"`
	Status    string    `bson:"status" json:"status"` // "ISSUED", "USED"
}

type ReceiptRepo interface {
	Save(receipt *Receipt) error
	GetByOrderID(orderID string) (*Receipt, error)
}
