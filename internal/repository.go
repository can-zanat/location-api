package internal

import (
	"context"
	"location-api/configs"
	"location-api/model"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	CreateLocation(req *model.CreateLocationRequest) (*model.CreateLocationResponse, error)
	GetLocation(req *model.GetLocationRequest) (*model.GetLocationResponse, error)
}

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

func (store *MongoDBStore) CreateLocation(req *model.CreateLocationRequest) (*model.CreateLocationResponse, error) {
	collection := store.Client.Database("location").Collection("locations")

	doc := bson.M{
		"name":         req.Name,
		"latitude":     req.Latitude,
		"longitude":    req.Longitude,
		"marker_color": req.MarkerColor,
		"created_at":   time.Now(),
	}

	result, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, mongo.ErrNilDocument
	}

	return &model.CreateLocationResponse{ID: insertedID.Hex()}, nil
}

func (store *MongoDBStore) GetLocation(req *model.GetLocationRequest) (*model.GetLocationResponse, error) {
	collection := store.Client.Database("location").Collection("locations")

	objectID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var location model.GetLocationResponse
	if err := collection.FindOne(context.TODO(), filter).Decode(&location); err != nil {
		return nil, err
	}

	if location.ID == "" {
		return nil, mongo.ErrNoDocuments
	}

	return &model.GetLocationResponse{
		ID:          location.ID,
		Name:        location.Name,
		Latitude:    location.Latitude,
		Longitude:   location.Longitude,
		MarkerColor: location.MarkerColor,
	}, nil
}
