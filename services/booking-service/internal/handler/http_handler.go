package handler

import (
	"github.com/dwikikusuma/ticket-rush/services/booking-service/internal/service"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	service *service.BookingService
}

func NewBookingHandler(service *service.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) RegisterRoutes(r gin.IRouter) {
	r.POST("/bookings", h.CreateBooking)
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("userId")

	var req struct {
		EventName string `json:"event_name"`
		Seat      string `json:"seat"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.CreateBooking(ctx, userID, req.EventName, req.Seat)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "booking created successfully"})
}
