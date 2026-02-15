package service

import "github.com/dwikikusuma/ticket-rush/services/fulfillment-service/internal/domain"

type Service struct {
	repo domain.ReceiptRepo
}

func NewService(repo domain.ReceiptRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveReceipt(receipt *domain.Receipt) error {
	return s.repo.Save(receipt)
}

func (s *Service) GetReceiptByOrderID(orderID string) (*domain.Receipt, error) {
	receipt, err := s.repo.GetByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
