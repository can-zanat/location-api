package internal

import (
	"context"
	"location-api/configs"
	"location-api/internal/helper"
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
	GetLocations(req *model.GetLocationsRequest) (*model.GetLocationsResponse, error)
	UpdateLocations(req *model.UpdateLocationsRequest) (*model.UpdateLocationsResponse, error)
	GetRoutes() (*model.GetAllLocationsDBResponse, error)
}

type MongoDBStore struct {
	Client *mongo.Client
}

const cacheKey = "cached_db_locations"
const cacheDuration = 30 * time.Second
const dbTimeout = 5 * time.Minute

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

func (store *MongoDBStore) GetLocations(req *model.GetLocationsRequest) (*model.GetLocationsResponse, error) {
	collection := store.Client.Database("location").Collection("locations")

	var page, limit = int64(req.Page), int64(req.Limit)
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	opts := options.Find().SetSkip(skip).SetLimit(limit)
	filter := bson.M{}

	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var locations []model.GetLocationResponse

	for cursor.Next(context.TODO()) {
		var location model.GetLocationResponse
		if err := cursor.Decode(&location); err != nil {
			return nil, err
		}

		locations = append(locations, location)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &model.GetLocationsResponse{Locations: locations}, nil
}

func (store *MongoDBStore) UpdateLocations(req *model.UpdateLocationsRequest) (*model.UpdateLocationsResponse, error) {
	collection := store.Client.Database("location").Collection("locations")

	var updatedIDs []string

	var failedIDs []string

	var totalModified int64

	for _, location := range req.Locations {
		objectID, err := primitive.ObjectIDFromHex(location.ID)
		if err != nil {
			failedIDs = append(failedIDs, location.ID)
			continue
		}

		updateData := bson.M{}
		orConditions := []bson.M{}

		if location.Name != "" {
			updateData["name"] = location.Name
			orConditions = append(orConditions, bson.M{"name": bson.M{"$ne": location.Name}})
		}

		if location.Latitude != 0.0 {
			updateData["latitude"] = location.Latitude
			orConditions = append(orConditions, bson.M{"latitude": bson.M{"$ne": location.Latitude}})
		}

		if location.Longitude != 0.0 {
			updateData["longitude"] = location.Longitude
			orConditions = append(orConditions, bson.M{"longitude": bson.M{"$ne": location.Longitude}})
		}

		if location.MarkerColor != "" {
			updateData["marker_color"] = location.MarkerColor
			orConditions = append(orConditions, bson.M{"marker_color": bson.M{"$ne": location.MarkerColor}})
		}

		if len(updateData) == 0 {
			failedIDs = append(failedIDs, location.ID)
			continue
		}

		filter := bson.M{
			"_id": objectID,
			"$or": orConditions,
		}

		updateData["updated_at"] = time.Now()
		update := bson.M{"$set": updateData}

		result, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			failedIDs = append(failedIDs, location.ID)
			continue
		}

		if result.ModifiedCount > 0 {
			updatedIDs = append(updatedIDs, location.ID)
			totalModified += result.ModifiedCount
		} else {
			failedIDs = append(failedIDs, location.ID)
		}
	}

	if len(updatedIDs) == 0 && len(failedIDs) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return &model.UpdateLocationsResponse{
		UpdatedIDs:   updatedIDs,
		FailedIDs:    failedIDs,
		UpdatedCount: totalModified,
	}, nil
}

func (store *MongoDBStore) GetRoutes() (*model.GetAllLocationsDBResponse, error) {
	var cachedLocations model.GetAllLocationsDBResponse
	if err := helper.GetCache(cacheKey, &cachedLocations); err == nil {
		log.Println("INFO: Data get from cache.")
		return &cachedLocations, nil
	}

	log.Println("WARNING: Redis empty, data will get from db...")

	collection := store.Client.Database("location").Collection("locations")

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var locations []model.GetLocationResponse

	for cursor.Next(ctx) {
		var location model.GetLocationResponse
		if err := cursor.Decode(&location); err != nil {
			return nil, err
		}

		locations = append(locations, location)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		log.Println("ERROR: No documents found in database.")
		return nil, mongo.ErrNoDocuments
	}

	dbResponse := &model.GetAllLocationsDBResponse{Locations: locations}
	_ = helper.SetCache(cacheKey, dbResponse, cacheDuration)

	log.Println("INFO: Data write redis.")

	return &model.GetAllLocationsDBResponse{Locations: locations}, nil
}
