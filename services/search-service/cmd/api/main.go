package main

import (
	"database/sql"
	"log"
	"net"
	"sync"
	"time"

	ticketV1 "github.com/dwikikusuma/ticket-rush/common/gen/ticket/v1"
	"github.com/dwikikusuma/ticket-rush/common/pkg/db"
	"github.com/dwikikusuma/ticket-rush/common/pkg/middleware"
	"github.com/dwikikusuma/ticket-rush/services/search-service/internal/adapter"
	ticketHandler "github.com/dwikikusuma/ticket-rush/services/search-service/internal/handler"
	"github.com/dwikikusuma/ticket-rush/services/search-service/internal/repository"
	ticketRPC "github.com/dwikikusuma/ticket-rush/services/search-service/internal/rpc"
	ticketSvc "github.com/dwikikusuma/ticket-rush/services/search-service/internal/service"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	elasticURL  = "http://localhost:9200"
	port        = "8083"
	pricingAddr = "50051"
	rpcAddr     = "50061"

	pHost     = "localhost"
	pUser     = "user"
	pPassword = "password"
	pDB       = "ticket_db"
	pPort     = 5432
)

func main() {

	es := openESConnection()
	ticketDB := openDBConnection()

	elasticRepo := repository.NewElasticRepo(es)
	pricingClient, err := adapter.NewPricingClient(pricingAddr)
	if err != nil {
		log.Fatalf("Error creating pricing client: %v", err)
	}

	dbRepo := repository.NewTicketRepo(ticketDB)
	service := ticketSvc.NewSearchService(elasticRepo, pricingClient, dbRepo)
	handler := ticketHandler.NewSearchHandler(service)

	grpcServer := grpc.NewServer()
	rpc := ticketRPC.NewTicketRpc(service)
	ticketV1.RegisterTicketServiceServer(grpcServer, rpc)
	reflection.Register(grpcServer)

	r := gin.Default()
	lis, err := net.Listen("tcp", ":"+rpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", rpcAddr, err)
	}

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

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Starting gRPC server on port %s", rpcAddr)

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	wg.Wait()
	log.Println("Search service stopped")
}

func openDBConnection() *sql.DB {
	log.Println("starting Database connection")
	dbConfig := db.Config{
		DB:              pDB,
		Pass:            pPassword,
		Port:            pPort,
		User:            pUser,
		Host:            pHost,
		MaxIdleConns:    3,
		MaxOpenConns:    10,
		ConnMaxLifetime: 1 * time.Hour,
	}

	dbConn, err := db.Open(dbConfig)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	return dbConn
}

func openESConnection() *elasticsearch.Client {
	log.Println("Starting search service")
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{elasticURL},
	})

	if err != nil {
		log.Fatalf("Error creating the client: %v", err)
	}
	log.Println("Created elastic search client")
	return es
}
