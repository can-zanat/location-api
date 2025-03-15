package internal

import (
	"context"
	"log"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
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

		req := &testLocationReq

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
