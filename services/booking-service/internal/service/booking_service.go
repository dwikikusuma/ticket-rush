package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	ticketPB "github.com/dwikikusuma/ticket-rush/common/gen/ticket/v1"
	"github.com/dwikikusuma/ticket-rush/services/booking-service/internal/domain"
	"github.com/redis/go-redis/v9"
)

type BookingService struct {
	repo        domain.BookingRepo
	ticketSvc   ticketPB.TicketServiceClient
	redisClient *redis.Client
}

func NewBookingService(repo domain.BookingRepo, ticketSvc ticketPB.TicketServiceClient, redisClient *redis.Client) *BookingService {
	return &BookingService{
		repo:        repo,
		ticketSvc:   ticketSvc,
		redisClient: redisClient,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, userID, eventName, seat string) error {
	lockKey := s.getTicketCacheKey(eventName, seat)

	if err := s.acquireLock(ctx, lockKey); err != nil {
		return err
	}
	defer s.releaseLock(ctx, lockKey)

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

	_, err = s.ticketSvc.UpdateTicketStatus(ctx, &ticketPB.UpdateTicketStatusRequest{
		TicketId: ticketDetail.TicketId,
		Status:   "Sold",
	})

	if err != nil {
		log.Printf("failed to update ticket status: %v", err)
		return err
	}

	err = s.repo.CreateBooking(ctx, domain.Booking{
		UserID:   userID,
		Status:   "Booked",
		TicketID: int(ticketDetail.TicketId),
	})

	if err != nil {
		log.Printf("failed to create booking: %v", err)
		_, _ = s.ticketSvc.UpdateTicketStatus(ctx, &ticketPB.UpdateTicketStatusRequest{
			TicketId: ticketDetail.TicketId,
			Status:   "AVAILABLE",
		})
		return err
	}

	return nil
}

func (s *BookingService) acquireLock(ctx context.Context, key string) error {
	acquired, err := s.redisClient.SetNX(ctx, key, "locked", 10*time.Second).Result()

	if err != nil {
		return fmt.Errorf("redis connection error: %v", err)
	}

	if !acquired {
		return errors.New("ticket is currently being booked by another user")
	}

	return nil
}

func (s *BookingService) releaseLock(ctx context.Context, key string) error {
	_, err := s.redisClient.Del(ctx, key).Result()
	return err
}

func (s *BookingService) getTicketCacheKey(eventName, seat string) string {
	return "lock:" + eventName + ":" + seat
}
