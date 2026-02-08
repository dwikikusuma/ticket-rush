package service

import (
	"context"
	"errors"
	"sync"

	"github.com/dwikikusuma/ticket-rush/services/search-service/internal/domain"
)

type searchService struct {
	repo          domain.SearchRepository
	dbRepo        domain.TicketRepo
	pricingClient domain.PricingClient
}

func NewSearchService(repo domain.SearchRepository, client domain.PricingClient, dbRepo domain.TicketRepo) domain.SearchService {
	return &searchService{
		repo:          repo,
		pricingClient: client,
		dbRepo:        dbRepo,
	}
}

func (s *searchService) FindTickets(query string, limit int, cursor string) (*domain.SearchResult, error) {
	result, err := s.repo.SearchQuery(query, limit, cursor)
	if err != nil {
		return &domain.SearchResult{}, err
	}

	if result == nil || len(result.Tickets) == 0 {
		return &domain.SearchResult{
			Tickets:    []domain.Ticket{},
			NextCursor: "",
		}, nil
	}

	var wg sync.WaitGroup
	for idx := range result.Tickets {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			t := result.Tickets[idx]
			pricing, _, err := s.pricingClient.GetRealTimePrice(context.Background(), &t)
			if err == nil {
				t.Price = int(pricing)
			}

		}(idx)
	}

	wg.Wait()
	return result, nil
}

func (s *searchService) FindTicketByID(ctx context.Context, id int) (*domain.Ticket, error) {
	return s.dbRepo.GetTicketByID(ctx, id)
}

func (s *searchService) FindTicketBySeatAndEvent(ctx context.Context, eventName string, seatID string) (*domain.Ticket, error) {
	return s.dbRepo.GetEventSeat(ctx, eventName, seatID)
}

func (s *searchService) UpdateTicketStatus(ctx context.Context, id int, status string) error {
	if status != "Sold" && status != "AVAILABLE" {
		return errors.New("invalid status")
	}
	return s.dbRepo.UpdateTicketStatus(ctx, id, status)
}
