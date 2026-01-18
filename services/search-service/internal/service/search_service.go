package service

import "github.com/dwikikusuma/ticket-rush/services/search-service/internal/domain"

type searchService struct {
	repo domain.SearchRepository
}

func NewSearchService(repo domain.SearchRepository) domain.SearchService {
	return &searchService{repo: repo}
}

func (s *searchService) FindTickets(query string, limit int, cursor string) (domain.SearchResult, error) {
	result, err := s.repo.SearchQuery(query, limit, cursor)
	if err != nil {
		return domain.SearchResult{}, err
	}
	return *result, nil
}
