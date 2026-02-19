package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	ticketV1 "github.com/dwikikusuma/ticket-rush/common/gen/ticket/v1"
	"github.com/dwikikusuma/ticket-rush/common/pkg/db"
	"github.com/dwikikusuma/ticket-rush/common/pkg/events"
	"github.com/dwikikusuma/ticket-rush/common/pkg/middleware"
	handler2 "github.com/dwikikusuma/ticket-rush/services/booking-service/internal/handler"
	"github.com/dwikikusuma/ticket-rush/services/booking-service/internal/repository"
	"github.com/dwikikusuma/ticket-rush/services/booking-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ticketAddr  = "localhost:50061"
	port        = "8086"
	kafkaBroker = "localhost:9092"
)

func main() {
	dbConn := openDBConnection()
	defer dbConn.Close()

	ticketConn, err := grpc.NewClient(ticketAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to ticket service: %v", err)
	}
	ticketSVC := ticketV1.NewTicketServiceClient(ticketConn)
	redisClient := newRedisClient()
	producer := events.NewProducer([]string{kafkaBroker})
	defer producer.Close()

	bookRepo := repository.NewBookingRepo(dbConn)
	bookSvc := service.NewBookingService(bookRepo, ticketSVC, redisClient, producer)
	handler := handler2.NewBookingHandler(bookSvc)

	r := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	r.Use(
		middleware.AuthMiddleware(),
		middleware.RequestID(),
		middleware.TimeOut(3*time.Second),
		gin.Logger(),
		gin.Recovery(),
	)
	handler.RegisterRoutes(r)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)

	}
}

func openDBConnection() *sql.DB {
	config := db.Config{
		Host: "localhost",
		User: "user",
		Pass: "password",
		DB:   "ticket_db",
		Port: 5432,
	}

	conn, err := db.Open(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return conn
}

func newRedisClient() *redis.Client {
	rc := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rc.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return rc
}
