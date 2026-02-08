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
