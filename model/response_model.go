package model

type CreateLocationResponse struct {
	ID string `json:"id" bson:"_id"`
}

type GetLocationResponse struct {
	ID          string  `json:"id" bson:"_id"`
	Name        string  `json:"name" bson:"name"`
	Latitude    float64 `json:"latitude" bson:"latitude"`
	Longitude   float64 `json:"longitude" bson:"longitude"`
	MarkerColor string  `json:"marker_color" bson:"marker_color"`
}

type GetLocationsResponse struct {
	Locations []GetLocationResponse `json:"locations"`
}
