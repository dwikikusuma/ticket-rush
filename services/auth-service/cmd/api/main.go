package main

import (
	"log"
	"sync"
	"time"

	"github.com/dwikikusuma/ticket-rush/common/pkg/db"
	"github.com/dwikikusuma/ticket-rush/services/auth-service/internal/handler"
	"github.com/dwikikusuma/ticket-rush/services/auth-service/internal/repository"
	"github.com/dwikikusuma/ticket-rush/services/auth-service/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	port      = ":8087"
	pHost     = "localhost"
	pUser     = "user"
	pPassword = "password"
	pDB       = "ticket_db"
	pPort     = 5432
)

func main() {
	log.Println("Starting Auth Service")

	postgresConfig := db.Config{
		DB:              pDB,
		Pass:            pPassword,
		Port:            pPort,
		User:            pUser,
		Host:            pHost,
		MaxIdleConns:    3,
		MaxOpenConns:    10,
		ConnMaxLifetime: 1 * time.Hour,
	}

	dbConn, err := db.Open(postgresConfig)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	userRepo := repository.NewUserRepo(dbConn)
	authSvc := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authSvc)

	r := gin.Default()
	authHandler.RegisterRoutes(r)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := r.Run(port); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()
	log.Println("Server started on port", port)
	wg.Wait()
	log.Println("Server stopped")
}
