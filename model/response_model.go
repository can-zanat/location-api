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

type UpdateLocationsResponse struct {
	UpdatedIDs   []string `json:"updated_ids"`
	FailedIDs    []string `json:"failed_ids"`
	UpdatedCount int64    `json:"updated_count"`
}

type Route struct {
	ID          string  `json:"id" bson:"_id"`
	Name        string  `json:"name" bson:"name"`
	Distance    float64 `json:"distance" bson:"distance"`
	MarkerColor string  `json:"marker_color" bson:"marker_color"`
}

type GetRoutesResponse struct {
	Routes []Route `json:"routes"`
}

type GetAllLocationsDBResponse struct {
	Locations []GetLocationResponse `json:"locations"`
}
