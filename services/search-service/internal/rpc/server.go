package rpc

import (
	"context"

	pb "github.com/dwikikusuma/ticket-rush/common/gen/ticket/v1"
	"github.com/dwikikusuma/ticket-rush/services/search-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TicketRpc struct {
	pb.UnimplementedTicketServiceServer
	svc domain.SearchService
}

func NewTicketRpc(svc domain.SearchService) *TicketRpc {
	return &TicketRpc{svc: svc}
}

func (r *TicketRpc) GetTicket(ctx context.Context, req *pb.TicketRequest) (*pb.TicketResponse, error) {
	ticket, err := r.svc.FindTicketByID(ctx, int(req.EventId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "ticket not found for ID %d: %v", req.EventId, err)
	}

	return &pb.TicketResponse{
		TicketId:  int32(ticket.ID),
		EventName: ticket.EventName,
		Stadium:   ticket.Stadium,
		Price:     int64(ticket.Price),
		SeatId:    ticket.SeatID,
		Status:    ticket.Status,
	}, nil
}
