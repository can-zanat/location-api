package model

type UpdateLocationsDBModel struct {
	Name        string  `json:"name" bson:"name"`
	Latitude    float64 `json:"latitude" bson:"latitude"`
	Longitude   float64 `json:"longitude" bson:"longitude"`
	MarkerColor string  `json:"marker_color" bson:"marker_color"`
}
