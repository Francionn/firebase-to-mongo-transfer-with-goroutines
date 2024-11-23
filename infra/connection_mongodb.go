package infra

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBHandler struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

const (
	databaseName   = "quality_air_data"
	collectionName = "air_data"
)

func NewMongoDBHandler() (*MongoDBHandler, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected MongoDB")

	database := client.Database(databaseName)
	collection := database.Collection(collectionName)

	return &MongoDBHandler{
		Client:     client,
		Database:   database,
		Collection: collection,
	}, nil
}

func (handler *MongoDBHandler) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := handler.Client.Disconnect(ctx); err != nil {
		log.Printf("Erro ao desconectar do MongoDB: %v", err)
	}
}
