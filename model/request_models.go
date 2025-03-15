package model

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type CreateLocationRequest struct {
	Name        string  `json:"name" bson:"name" validate:"required,min=3"`
	Latitude    float64 `json:"latitude" bson:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" bson:"longitude" validate:"required"`
	MarkerColor string  `json:"marker_color" bson:"marker_color" validate:"required,len=6,hexadecimal"`
}

type GetLocationRequest struct {
	ID string `query:"id" json:"id" bson:"_id" validate:"required"`
}

type GetLocationsRequest struct {
	Page  int `query:"page" json:"page" bson:"page" validate:"required"`
	Limit int `query:"limit" json:"limit" bson:"limit" validate:"required"`
}

type UpdateLocation struct {
	ID          string  `json:"id" bson:"_id" validate:"required"`
	Name        string  `json:"name" bson:"name" validate:"omitempty,min=3"`
	Latitude    float64 `json:"latitude" bson:"latitude" validate:"omitempty"`
	Longitude   float64 `json:"longitude" bson:"longitude" validate:"omitempty"`
	MarkerColor string  `json:"marker_color" bson:"marker_color" validate:"omitempty,len=6,hexadecimal"`
}

type UpdateLocationsRequest struct {
	Locations []UpdateLocation `json:"locations" bson:"locations" validate:"required,dive"`
}

func (req *CreateLocationRequest) ValidateLocation() error {
	return validate.Struct(req)
}
