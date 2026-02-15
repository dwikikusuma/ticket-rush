package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dwikikusuma/ticket-rush/common/pkg/events"
	"github.com/dwikikusuma/ticket-rush/services/fulfillment-service/internal/domain"
	"github.com/dwikikusuma/ticket-rush/services/fulfillment-service/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	kafkaBroker   = "localhost:9092"
	consumerGroup = "fulfillment-service-consumer-group-v2"
	bookingTopic  = "booking.created"
	mongoURI      = "mongodb://user:password@localhost:27017"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := setupMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("ticket_rush").Collection("receipts")
	repo := repository.NewMongoRepo(collection)

	consumer := events.NewConsumer([]string{kafkaBroker}, consumerGroup, bookingTopic)
	defer consumer.Close()

	consumerWorker := domain.NewConsumerWorker(consumer, repo)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumerWorker.Start(ctx)
	}()

	log.Println("Starting consumer...")
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-sigChannel

	log.Println("Shutdown signal received, stopping consumer...")
	cancel()
	wg.Wait()
	log.Println("Consumer stopped gracefully")
}

func setupMongoDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
