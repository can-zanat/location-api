package internal

import (
	"context"
	"location-api/model"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoImage = "mongo:7.0.4"
)

func NewStoreWithURI(uri string) *MongoDBStore {
	clientOptions := options.Client().ApplyURI(uri)
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

func prepareTestStore(t *testing.T) (store *MongoDBStore, clean func()) {
	t.Helper()

	ctx := context.Background()

	mongodbContainer, err := mongodb.RunContainer(ctx, testcontainers.WithImage(mongoImage))
	if err != nil {
		t.Fatalf("Failed to start MongoDB container: %v", err)
	}

	clean = func() {
		if terminateErr := mongodbContainer.Terminate(ctx); terminateErr != nil {
			t.Fatalf("Failed to terminate MongoDB container: %v", terminateErr)
		}
	}

	containerURI, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Failed to get container connection string: %v", err)
	}

	s := NewStoreWithURI(containerURI)

	return s, clean
}

func TestMongoDBStore_CreateLocation(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Run("should create a location", func(t *testing.T) {
		store, clean := prepareTestStore(t)
		defer clean()

		req := &testCreateLocationReq

		resp, err := store.CreateLocation(req)
		if err != nil {
			t.Fatalf("Failed to create location: %v", err)
		}

		if resp.ID == "" {
			t.Fatalf("Expected location ID to be set, got empty string")
		}

		t.Logf("Created location with ID: %s", resp.ID)
	})
}

func TestMongoDBStore_GetLocation(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Run("should return error document not found", func(t *testing.T) {
		store, clean := prepareTestStore(t)
		defer clean()

		insertedID := "5f9b1f3b1c9d440000f1b4b0"

		req := &model.GetLocationRequest{
			ID: insertedID,
		}

		resp, err := store.GetLocation(req)

		if err == nil {
			t.Fatalf("Expected error, but got nil")
		}

		assert.Nil(t, resp)
		assert.Equal(t, mongo.ErrNoDocuments, err)
	})

	t.Run("should get a location", func(t *testing.T) {
		store, clean := prepareTestStore(t)
		defer clean()

		collection := store.Client.Database("location").Collection("locations")
		locationDoc := bson.M{
			"name":         "test3",
			"latitude":     124.12,
			"longitude":    134.12,
			"marker_color": "FFFAFF",
		}

		result, err := collection.InsertOne(context.Background(), locationDoc)
		if err != nil {
			t.Fatalf("Failed to insert location: %v", err)
		}

		insertedID, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			t.Fatalf("Failed to convert inserted ID to ObjectID")
		}

		req := &model.GetLocationRequest{
			ID: insertedID.Hex(),
		}

		resp, err := store.GetLocation(req)
		if err != nil {
			t.Fatalf("Failed to get location: %v", err)
		}

		if resp.ID == "" {
			t.Fatalf("Expected location ID to be set, got empty string")
		}

		t.Logf("Got location with ID: %s", resp.ID)
	})
}

func TestMongoDBStore_GetLocations(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Run("should return no locations found", func(t *testing.T) {
		store, clean := prepareTestStore(t)
		defer clean()

		req := &model.GetLocationsRequest{}

		resp, err := store.GetLocations(req)
		if err == nil {
			t.Fatalf("Expected error, but got nil")
		}

		assert.Nil(t, resp)
		assert.Equal(t, mongo.ErrNoDocuments, err)
	})
	t.Run("should get locations", func(t *testing.T) {
		store, clean := prepareTestStore(t)
		defer clean()

		collection := store.Client.Database("location").Collection("locations")
		locationDoc := bson.M{
			"name":         "test3",
			"latitude":     124.12,
			"longitude":    134.12,
			"marker_color": "FFFAFF",
		}

		_, err := collection.InsertOne(context.Background(), locationDoc)
		if err != nil {
			t.Fatalf("Failed to insert location: %v", err)
		}

		req := &model.GetLocationsRequest{}

		resp, err := store.GetLocations(req)
		if err != nil {
			t.Fatalf("Failed to get locations: %v", err)
		}

		if len(resp.Locations) == 0 {
			t.Fatalf("Expected locations to be returned, got empty slice")
		}

		t.Logf("Got %d locations", len(resp.Locations))
	})
}
