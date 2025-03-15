package internal

import "location-api/model"

type Service struct {
	store Store
}

type LocationDBStore interface {
	CreateLocation(req *model.CreateLocationRequest) (*model.CreateLocationResponse, error)
	GetLocation(req *model.GetLocationRequest) (*model.GetLocationResponse, error)
}

func NewService(s Store) *Service {
	return &Service{store: s}
}

func (s *Service) CreateLocation(req *model.CreateLocationRequest) (*model.CreateLocationResponse, error) {
	return s.store.CreateLocation(req)
}

func (s *Service) GetLocation(req *model.GetLocationRequest) (*model.GetLocationResponse, error) {
	return s.store.GetLocation(req)
}
