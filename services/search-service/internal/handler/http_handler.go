package handler

import (
	"net/http"
	"strconv"

	"github.com/dwikikusuma/ticket-rush/services/search-service/internal/domain"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	service domain.SearchService
}

func NewSearchHandler(service domain.SearchService) *SearchHandler {
	return &SearchHandler{
		service: service,
	}
}

func (h *SearchHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/search", h.Search)
}

func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	limitStr := c.DefaultQuery("limit", "10")
	cursor := c.DefaultQuery("cursor", "")

	limit, _ := strconv.Atoi(limitStr)
	result, err := h.service.FindTickets(query, limit, cursor)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":        result.Tickets,
		"next_cursor": result.NextCursor,
	})
}
