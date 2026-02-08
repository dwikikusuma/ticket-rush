package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/dwikikusuma/ticket-rush/services/booking-service/internal/domain"
	bookingDB "github.com/dwikikusuma/ticket-rush/services/booking-service/internal/infra/postgres"
	"github.com/google/uuid"
)

type bookingRepo struct {
	db *bookingDB.Queries
}

func NewBookingRepo(db *sql.DB) domain.BookingRepo {
	return &bookingRepo{db: bookingDB.New(db)}
}

func (r *bookingRepo) CreateBooking(ctx context.Context, req domain.Booking) error {
	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Printf("invalid user ID: %v", err)
		return err
	}
	params := bookingDB.CreateBookingParams{
		UserID:   userUUID,
		TicketID: int32(req.TicketID),
		Status:   req.Status,
	}
	_, err = r.db.CreateBooking(ctx, params)

	return err
}

func (r *bookingRepo) UpdateBookingStatusFailed(ctx context.Context, bookingID string) error {
	bookingUUID, err := uuid.Parse(bookingID)
	if err != nil {
		log.Printf("invalid booking ID: %v", err)
		return err
	}

	err = r.db.UpdateBookingStatus(ctx, bookingDB.UpdateBookingStatusParams{
		Status: "Failed",
		ID:     bookingUUID,
	})

	if err != nil {
		log.Printf("failed to update booking status: %v", err)
		return err
	}

	return nil
}
