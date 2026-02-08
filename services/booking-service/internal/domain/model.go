package domain

import "time"

type Booking struct {
	ID        string
	UserID    string
	TicketID  int
	Status    string
	CreatedAt time.Time
}
