package main

import (
	"database/sql"
	"log"
	"time"

	ticketV1 "github.com/dwikikusuma/ticket-rush/common/gen/ticket/v1"
	"github.com/dwikikusuma/ticket-rush/common/pkg/db"
	"github.com/dwikikusuma/ticket-rush/common/pkg/middleware"
	handler2 "github.com/dwikikusuma/ticket-rush/services/booking-service/internal/handler"
	"github.com/dwikikusuma/ticket-rush/services/booking-service/internal/repository"
	"github.com/dwikikusuma/ticket-rush/services/booking-service/internal/service"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ticketAddr = "localhost:50061"
	port       = "8085"
)

func main() {
	dbConn := openDBConnection()
	defer dbConn.Close()

	ticketConn, err := grpc.NewClient(ticketAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to ticket service: %v", err)
	}
	ticketSVC := ticketV1.NewTicketServiceClient(ticketConn)

	bookRepo := repository.NewBookingRepo(dbConn)
	bookSvc := service.NewBookingService(bookRepo, ticketSVC)
	handler := handler2.NewBookingHandler(bookSvc)

	r := gin.Default()
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
