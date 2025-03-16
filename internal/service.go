package internal

import (
	"location-api/internal/helper"
	"location-api/model"
	"sort"
)

type Service struct {
	store Store
}

type LocationDBStore interface {
	CreateLocation(req *model.CreateLocationRequest) (*model.CreateLocationResponse, error)
	GetLocation(req *model.GetLocationRequest) (*model.GetLocationResponse, error)
	GetLocations(req *model.GetLocationsRequest) (*model.GetLocationsResponse, error)
	UpdateLocations(req *model.UpdateLocationsRequest) (*model.UpdateLocationsResponse, error)
	GetRoutes() (*model.GetAllLocationsDBResponse, error)
}

func NewService(s Store) *Service {
	return &Service{store: s}
}

func (s *Service) CreateLocation(req *model.CreateLocationRequest) (*model.CreateLocationResponse, error) {
	res, err := s.store.CreateLocation(req)
	if err != nil {
		return nil, err
	}

	_ = helper.DeleteCache("cached_db_locations")

	return res, nil
}

func (s *Service) GetLocation(req *model.GetLocationRequest) (*model.GetLocationResponse, error) {
	return s.store.GetLocation(req)
}

func (s *Service) GetLocations(req *model.GetLocationsRequest) (*model.GetLocationsResponse, error) {
	return s.store.GetLocations(req)
}

func (s *Service) UpdateLocations(req *model.UpdateLocationsRequest) (*model.UpdateLocationsResponse, error) {
	res, err := s.store.UpdateLocations(req)
	if err != nil {
		return nil, err
	}

	_ = helper.DeleteCache("cached_db_locations")

	return res, nil
}

func (s *Service) GetRoutes(req *model.GetRoutesRequest) (*model.GetRoutesResponse, error) {
	locationsResp, err := s.store.GetRoutes()
	if err != nil {
		return nil, err
	}

	if len(locationsResp.Locations) == 0 {
		return &model.GetRoutesResponse{Routes: []model.Route{}}, nil
	}

	if req.Latitude == 0 || req.Longitude == 0 {
		return nil, err
	}

	type LocationDistance struct {
		Location model.GetLocationResponse
		Distance float64
	}

	locationDistances := make([]LocationDistance, 0, len(locationsResp.Locations))

	for _, loc := range locationsResp.Locations {
		distance := helper.Haversine(req.Latitude, req.Longitude, loc.Latitude, loc.Longitude)
		locationDistances = append(locationDistances, LocationDistance{
			Location: loc,
			Distance: distance,
		})
	}

	sort.Slice(locationDistances, func(i, j int) bool {
		return locationDistances[i].Distance < locationDistances[j].Distance
	})

	sortedRoutes := make([]model.Route, 0, len(locationDistances))
	for _, loc := range locationDistances {
		sortedRoutes = append(sortedRoutes, model.Route{
			ID:          loc.Location.ID,
			Name:        loc.Location.Name,
			Distance:    loc.Distance,
			MarkerColor: loc.Location.MarkerColor,
		})
	}

	return &model.GetRoutesResponse{
		Routes: sortedRoutes,
	}, nil
}
