package main

import (
	"log"
	"sync"
	"time"

	"github.com/dwikikusuma/ticket-rush/common/pkg/middleware"
	"github.com/dwikikusuma/ticket-rush/services/search-service/internal/adapter"
	ticketHandler "github.com/dwikikusuma/ticket-rush/services/search-service/internal/handler"
	"github.com/dwikikusuma/ticket-rush/services/search-service/internal/repository"
	ticketSvc "github.com/dwikikusuma/ticket-rush/services/search-service/internal/service"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

const (
	elasticURL  = "http://localhost:9200"
	port        = "8081"
	pricingAddr = "50051"
)

func main() {
	log.Println("Starting search service")
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{elasticURL},
	})

	log.Println("Created elastic search client")
	if err != nil {
		log.Fatalf("Error creating the client: %v", err)
	}

	elasticRepo := repository.NewElasticRepo(es)
	pricingClient, err := adapter.NewPricingClient(pricingAddr)
	if err != nil {
		log.Fatalf("Error creating pricing client: %v", err)
	}

	service := ticketSvc.NewSearchService(elasticRepo, pricingClient)
	handler := ticketHandler.NewSearchHandler(service)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "Search Service is Alive"})
	})

	protected := r.Group("/")
	protected.Use(
		middleware.AuthMiddleware(),
		middleware.RequestID(),
		middleware.TimeOut(3*time.Second),
		gin.Logger(),
		gin.Recovery(),
	)
	handler.RegisterRoutes(protected)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()
	wg.Wait()
	log.Println("Search service stopped")
}
