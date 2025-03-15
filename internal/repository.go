package internal

import (
	"context"
	"location-api/configs"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface{}

type MongoDBStore struct {
	Client *mongo.Client
}

func NewStore() *MongoDBStore {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	clientOptions := options.Client().ApplyURI(config.MongoDB.URI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Connection failure:", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatal("Unable to access MongoDB:", err)
	}

	return &MongoDBStore{
		Client: client,
	}
}
