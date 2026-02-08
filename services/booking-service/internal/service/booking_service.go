package service

import (
	"context"
	"errors"
	"log"

	ticketPB "github.com/dwikikusuma/ticket-rush/common/gen/ticket/v1"
	"github.com/dwikikusuma/ticket-rush/services/booking-service/internal/domain"
)

type BookingService struct {
	repo      domain.BookingRepo
	ticketSvc ticketPB.TicketServiceClient
}

func NewBookingService(repo domain.BookingRepo, ticketSvc ticketPB.TicketServiceClient) *BookingService {
	return &BookingService{
		repo:      repo,
		ticketSvc: ticketSvc,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, userID, eventName, seat string) error {

	ticketDetail, err := s.ticketSvc.GetTicketBySeatAndEvent(ctx, &ticketPB.TicketBySeatAndEventRequest{
		EventName: eventName,
		SeatId:    seat,
	})

	if err != nil {
		log.Printf("failed to get ticket details: %v", err)
		return err
	}

	if ticketDetail.Status != "AVAILABLE" {
		log.Printf("ticket is not available: %s", ticketDetail.Status)
		return errors.New("ticket is not available")
	}

	err = s.repo.CreateBooking(ctx, domain.Booking{
		UserID:   userID,
		Status:   "Booked",
		TicketID: int(ticketDetail.TicketId),
	})

	if err != nil {
		log.Printf("failed to create booking: %v", err)
		return err
	}

	_, err = s.ticketSvc.UpdateTicketStatus(ctx, &ticketPB.UpdateTicketStatusRequest{
		TicketId: ticketDetail.TicketId,
		Status:   "Sold",
	})

	if err != nil {
		log.Printf("failed to update ticket status: %v", err)
		_ = s.repo.UpdateBookingStatusFailed(ctx, userID)
		return err
	}

	return nil
}
