package domain

import "context"

type BookingRepo interface {
	CreateBooking(ctx context.Context, req Booking) error
}
