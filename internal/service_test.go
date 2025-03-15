package internal

import (
	"location-api/model"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var testLocationReq = model.CreateLocationRequest{
	Name:        "test",
	Latitude:    1.1,
	Longitude:   1.1,
	MarkerColor: "FFFFFF",
}

var testLocationRes = model.CreateLocationResponse{
	ID: "test",
}

func TestService_GetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("should create location properly", func(t *testing.T) {
		mockRepository := NewMockStore(ctrl)

		mockRepository.
			EXPECT().
			CreateLocation(&testLocationReq).
			Return(&testLocationRes, nil).
			Times(1)

		service := NewService(mockRepository)

		locationRes, _ := service.CreateLocation(&testLocationReq)
		assert.Equal(t, &testLocationRes, locationRes)
	})

	t.Run("return error", func(t *testing.T) {
		mockRepository := NewMockStore(ctrl)

		var err error

		expectedError := fiber.NewError(fiber.StatusInternalServerError, "Internal Server Error")

		mockRepository.
			EXPECT().
			CreateLocation(&testLocationReq).
			Return(nil, &fiber.Error{Code: 500, Message: "Internal Server Error"}).
			Times(1)

		service := NewService(mockRepository)

		_, err = service.CreateLocation(&testLocationReq)
		assert.Equal(t, expectedError, err)
	})
}
