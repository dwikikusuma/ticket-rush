package service

import (
	"context"
	"sync"

	"github.com/dwikikusuma/ticket-rush/services/search-service/internal/domain"
)

type searchService struct {
	repo          domain.SearchRepository
	pricingClient domain.PricingClient
}

func NewSearchService(repo domain.SearchRepository, client domain.PricingClient) domain.SearchService {
	return &searchService{
		repo:          repo,
		pricingClient: client,
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
