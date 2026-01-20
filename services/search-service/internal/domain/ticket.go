package domain

import "context"

type Ticket struct {
	ID        int    `json:"id"`
	EventName string `json:"event_name"`
	Stadium   string `json:"stadium"`
	Price     int    `json:"price"`
	SeatID    string `json:"seat_id"`
	Status    string `json:"status"`
}

type SearchResult struct {
	Tickets    []Ticket
	NextCursor string
}

type SearchRepository interface {
	SearchQuery(query string, limit int, cursor string) (*SearchResult, error)
}

type SearchService interface {
	FindTickets(query string, limit int, cursor string) (*SearchResult, error)
}

type PricingClient interface {
	GetRealTimePrice(ctx context.Context, ticket *Ticket) (int32, float32, error)
}
