package repository

import (
	"context"
	"database/sql"

	"github.com/dwikikusuma/ticket-rush/services/search-service/internal/domain"
	ticketDB "github.com/dwikikusuma/ticket-rush/services/search-service/internal/infra/postgres"
)

type ticketRepo struct {
	db *ticketDB.Queries
}

func NewTicketRepo(db *sql.DB) domain.TicketRepo {
	return &ticketRepo{db: ticketDB.New(db)}
}

func (r *ticketRepo) GetTicketByID(ctx context.Context, id int) (*domain.Ticket, error) {
	ticket, err := r.db.GetTicketByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return &domain.Ticket{
		ID:        int(ticket.ID),
		EventName: ticket.EventName.String,
		SeatID:    ticket.SeatID.String,
		Price:     int(ticket.Price.Int32),
		Status:    ticket.Status.String,
	}, nil
}

func (r *ticketRepo) GetEventSeat(ctx context.Context, eventName string, seatID string) (*domain.Ticket, error) {
	seatIDNS := sql.NullString{String: seatID, Valid: true}
	eventNameNS := sql.NullString{String: eventName, Valid: true}
	params := ticketDB.GetTicketBySeatAndEventParams{
		SeatID:    seatIDNS,
		EventName: eventNameNS,
	}

	ticket, err := r.db.GetTicketBySeatAndEvent(ctx, params)
	if err != nil {
		return nil, err
	}

	return &domain.Ticket{
		ID:        int(ticket.ID),
		EventName: ticket.EventName.String,
		SeatID:    ticket.SeatID.String,
		Price:     int(ticket.Price.Int32),
		Status:    ticket.Status.String,
	}, nil
}

func (r *ticketRepo) UpdateTicketStatus(ctx context.Context, id int, status string) error {
	params := ticketDB.UpdateTicketStatusParams{
		ID:     int32(id),
		Status: sql.NullString{String: status, Valid: true},
	}
	return r.db.UpdateTicketStatus(ctx, params)
}
