package repository

import (
	"context"
	"log"
	"time"

	"github.com/dwikikusuma/ticket-rush/services/fulfillment-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	collection *mongo.Collection
}

func NewMongoRepo(collection *mongo.Collection) domain.ReceiptRepo {
	return &mongoRepo{collection: collection}
}

func (r *mongoRepo) Save(receipt *domain.Receipt) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, receipt)
	if err != nil {
		log.Printf("failed to insert receipt: %v", err)
		return err
	}
	return nil
}

func (r *mongoRepo) GetByOrderID(orderID string) (*domain.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var receipt domain.Receipt
	err := r.collection.FindOne(ctx, bson.M{"booking_id": orderID}).Decode(&receipt)
	if err != nil {
		log.Printf("failed to find receipt: %v", err)
		return nil, err
	}

	return &receipt, nil
}
