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
	ticket, err := r.svc.FindTicketByID(ctx, int(req.TicketId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "ticket not found for ID %d: %v", req.TicketId, err)
	}

	return r.ticketToProto(*ticket), nil
}

func (r *TicketRpc) GetTicketBySeatAndEvent(ctx context.Context, req *pb.TicketBySeatAndEventRequest) (*pb.TicketResponse, error) {
	ticket, err := r.svc.FindTicketBySeatAndEvent(ctx, req.EventName, req.SeatId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "ticket not found for event %s and seat %s: %v", req.EventName, req.SeatId, err)
	}

	return r.ticketToProto(*ticket), nil
}

func (r *TicketRpc) UpdateTicketStatus(ctx context.Context, req *pb.UpdateTicketStatusRequest) (*pb.UpdateTicketStatusResponse, error) {
	if req.Status == "Sold" || req.Status == "AVAILABLE" {
		err := r.svc.UpdateTicketStatus(ctx, int(req.TicketId), req.Status)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update ticket status to Sold for ID %d: %v", req.TicketId, err)
		}
		return &pb.UpdateTicketStatusResponse{
			Status:   req.Status,
			TicketId: req.TicketId,
		}, nil
	} else {
		return nil, status.Error(codes.InvalidArgument, "invalid status")
	}
}

func (r *TicketRpc) ticketToProto(t domain.Ticket) *pb.TicketResponse {
	return &pb.TicketResponse{
		TicketId:  int32(t.ID),
		EventName: t.EventName,
		Stadium:   t.Stadium,
		Price:     int64(t.Price),
		SeatId:    t.SeatID,
		Status:    t.Status,
	}
}
