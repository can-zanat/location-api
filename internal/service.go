package internal

import "location-api/model"

type Service struct {
	store Store
}

type LocationDBStore interface {
	CreateLocation(req *model.CreateLocationRequest) (*model.CreateLocationResponse, error)
	GetLocation(req *model.GetLocationRequest) (*model.GetLocationResponse, error)
	GetLocations(req *model.GetLocationsRequest) (*model.GetLocationsResponse, error)
	UpdateLocations(req *model.UpdateLocationsRequest) (*model.UpdateLocationsResponse, error)
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

func (s *Service) GetLocations(req *model.GetLocationsRequest) (*model.GetLocationsResponse, error) {
	return s.store.GetLocations(req)
}

func (s *Service) UpdateLocations(req *model.UpdateLocationsRequest) (*model.UpdateLocationsResponse, error) {
	return s.store.UpdateLocations(req)
}
