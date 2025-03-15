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

func (req *CreateLocationRequest) ValidateLocation() error {
	return validate.Struct(req)
}
