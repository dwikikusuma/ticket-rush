package main

import (
	"log"
	"sync"

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
	handler.RegisterRoutes(r)

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
