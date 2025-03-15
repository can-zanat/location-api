package internal

import (
	"location-api/model"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var testCreateLocationReq = model.CreateLocationRequest{
	Name:        "test",
	Latitude:    1.1,
	Longitude:   1.1,
	MarkerColor: "FFFFFF",
}

var testCreateLocationRes = model.CreateLocationResponse{
	ID: "test",
}

var testGetLocationReq = model.GetLocationRequest{
	ID: "67d562e3d9f2d225ca4d9918",
}

var testGetLocationRes = model.GetLocationResponse{
	ID:          "test",
	Name:        "test",
	Latitude:    1.1,
	Longitude:   1.1,
	MarkerColor: "FFFFFF",
}

var testGetLocationsReq = model.GetLocationsRequest{
	Page:  1,
	Limit: 1,
}

var testGetLocationsRes = model.GetLocationsResponse{
	Locations: []model.GetLocationResponse{
		{
			ID:          "67d562e3d9f2d225ca4d9918",
			Name:        "test",
			Latitude:    1.1,
			Longitude:   1.1,
			MarkerColor: "FFFFFF",
		},
	},
}

var testUpdateLocationsReq = model.UpdateLocationsRequest{
	Locations: []model.UpdateLocation{
		{
			ID:          "67d562e3d9f2d225ca4d9918",
			Name:        "test",
			Latitude:    1.1,
			Longitude:   1.1,
			MarkerColor: "FFFFFF",
		},
		{
			ID:          "67d562e3d9f2d225ca4d9919",
			Name:        "test2",
			Latitude:    2.2,
			Longitude:   2.2,
			MarkerColor: "000000",
		},
	},
}

var testUpdateLocationsRes = model.UpdateLocationsResponse{
	UpdatedIDs:   []string{"67d562e3d9f2d225ca4d9918", "67d562e3d9f2d225ca4d9919"},
	FailedIDs:    []string{},
	UpdatedCount: 2,
}

func TestService_CreateLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("should create location properly", func(t *testing.T) {
		mockRepository := NewMockStore(ctrl)

		mockRepository.
			EXPECT().
			CreateLocation(&testCreateLocationReq).
			Return(&testCreateLocationRes, nil).
			Times(1)

		service := NewService(mockRepository)

		locationRes, _ := service.CreateLocation(&testCreateLocationReq)
		assert.Equal(t, &testCreateLocationRes, locationRes)
	})

	t.Run("return error", func(t *testing.T) {
		mockRepository := NewMockStore(ctrl)

		var err error

		expectedError := fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")

		mockRepository.
			EXPECT().
			CreateLocation(&testCreateLocationReq).
			Return(nil, &fiber.Error{Code: 500, Message: "Internal Server Error"}).
			Times(1)

		service := NewService(mockRepository)

		_, err = service.CreateLocation(&testCreateLocationReq)
		assert.Equal(t, expectedError, err)
	})
}

func TestService_GetLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("should get location properly", func(t *testing.T) {
		mockRepository := NewMockStore(ctrl)

		mockRepository.
			EXPECT().
			GetLocation(&testGetLocationReq).
			Return(&testGetLocationRes, nil).
			Times(1)

		service := NewService(mockRepository)

		locationRes, _ := service.GetLocation(&testGetLocationReq)
		assert.Equal(t, &testGetLocationRes, locationRes)
	})

	t.Run("return error", func(t *testing.T) {
		mockRepository := NewMockStore(ctrl)

		var err error

		expectedError := fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")

		mockRepository.
			EXPECT().
			GetLocation(&testGetLocationReq).
			Return(nil, &fiber.Error{Code: 500, Message: "Internal Server Error"}).
			Times(1)

		service := NewService(mockRepository)

		_, err = service.GetLocation(&testGetLocationReq)
		assert.Equal(t, expectedError, err)
	})
}

func TestService_GetLocations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("should get locations properly", func(t *testing.T) {
		mockRepository := NewMockStore(ctrl)

		mockRepository.
			EXPECT().
			GetLocations(&testGetLocationsReq).
			Return(&testGetLocationsRes, nil).
			Times(1)

		service := NewService(mockRepository)

		locationsRes, _ := service.GetLocations(&testGetLocationsReq)
		assert.Equal(t, &testGetLocationsRes, locationsRes)
	})

	t.Run("return error", func(t *testing.T) {
		mockRepository := NewMockStore(ctrl)

		var err error

		expectedError := fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")

		mockRepository.
			EXPECT().
			GetLocations(&testGetLocationsReq).
			Return(nil, &fiber.Error{Code: 500, Message: "Internal Server Error"}).
			Times(1)

		service := NewService(mockRepository)

		_, err = service.GetLocations(&testGetLocationsReq)
		assert.Equal(t, expectedError, err)
	})
}

func TestService_UpdateLocations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("should update locations properly", func(t *testing.T) {
		mockRepository := NewMockStore(ctrl)

		mockRepository.
			EXPECT().
			UpdateLocations(&testUpdateLocationsReq).
			Return(&testUpdateLocationsRes, nil).
			Times(1)

		service := NewService(mockRepository)

		locationsRes, err := service.UpdateLocations(&testUpdateLocationsReq)
		assert.Nil(t, err)
		assert.Equal(t, &testUpdateLocationsRes, locationsRes)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		mockRepository := NewMockStore(ctrl)

		expectedError := fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")

		mockRepository.
			EXPECT().
			UpdateLocations(&testUpdateLocationsReq).
			Return(nil, expectedError).
			Times(1)

		service := NewService(mockRepository)

		_, err := service.UpdateLocations(&testUpdateLocationsReq)
		assert.Equal(t, expectedError, err)
	})
}
