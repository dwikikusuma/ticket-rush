package domain

import "context"

type BookingRepo interface {
	CreateBooking(ctx context.Context, req Booking) error
	UpdateBookingStatusFailed(ctx context.Context, bookingID string) error
}
