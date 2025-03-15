package internal

import "location-api/model"

type Service struct {
	store Store
}

type LocationDBStore interface {
	CreateLocation(req *model.CreateLocationRequest) (*model.CreateLocationResponse, error)
}

func NewService(s Store) *Service {
	return &Service{store: s}
}

func (s *Service) CreateLocation(req *model.CreateLocationRequest) (*model.CreateLocationResponse, error) {
	return s.store.CreateLocation(req)
}
