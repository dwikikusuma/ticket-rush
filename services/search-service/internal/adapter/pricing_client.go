package adapter

import (
	"context"
	"time"

	pricingv1 "github.com/dwikikusuma/ticket-rush/common/gen/pricing/v1"
	"github.com/dwikikusuma/ticket-rush/services/search-service/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcPricingClient struct {
	client pricingv1.PricingServiceClient
}

func NewPricingClient(addr string) (domain.PricingClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pricingv1.NewPricingServiceClient(conn)
	return &grpcPricingClient{client: client}, nil
}

func (c *grpcPricingClient) GetRealTimePrice(ctx context.Context, ticket *domain.Ticket) (int32, float32, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	req := &pricingv1.PriceRequest{
		SeatId:    ticket.SeatID,
		EventId:   ticket.EventName,
		BasePrice: float32(ticket.Price),
	}

	price, err := c.client.GetPrice(ctx, req)
	if err != nil {
		return 0, 0, err
	}

	return price.GetFinalPrice(), price.GetMultiplier(), nil
}
