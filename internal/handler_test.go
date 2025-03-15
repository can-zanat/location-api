package internal

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/gofiber/fiber/v2"
)

func TestHandler_GetUsers(t *testing.T) {
	t.Run("should create location properly", func(t *testing.T) {
		mockService, mockServiceController := createMockService(t)
		defer mockServiceController.Finish()

		app := createTestApp()

		mockService.
			EXPECT().
			CreateLocation(&testLocationReq).
			Return(&testLocationRes, nil).
			Times(1)

		handler := NewHandler(mockService)
		handler.RegisterRoutes(app)

		req := httptest.NewRequest(
			http.MethodPost,
			"/location",
			bytes.NewReader([]byte(`{"name": "test", "latitude": 1.1, "longitude": 1.1, "marker_color": "FFFFFF"}`)),
		)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		defer res.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should return internal server error", func(t *testing.T) {
		mockService, mockServiceController := createMockService(t)
		defer mockServiceController.Finish()

		app := createTestApp()

		mockService.
			EXPECT().
			CreateLocation(&testLocationReq).
			Return(nil, assert.AnError).
			Times(1)

		handler := NewHandler(mockService)
		handler.RegisterRoutes(app)

		req := httptest.NewRequest(
			http.MethodPost,
			"/location",
			bytes.NewReader([]byte(`{"name": "test", "latitude": 1.1, "longitude": 1.1, "marker_color": "FFFFFF"}`)),
		)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		defer res.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("should return bad request error", func(t *testing.T) {
		mockService, mockServiceController := createMockService(t)
		defer mockServiceController.Finish()

		app := createTestApp()

		handler := NewHandler(mockService)
		handler.RegisterRoutes(app)

		req := httptest.NewRequest(
			http.MethodPost,
			"/location",
			http.NoBody,
		)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		defer res.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
	t.Run("should return validation error", func(t *testing.T) {
		mockService, mockServiceController := createMockService(t)
		defer mockServiceController.Finish()

		app := createTestApp()

		handler := NewHandler(mockService)
		handler.RegisterRoutes(app)

		req := httptest.NewRequest(
			http.MethodPost,
			"/location",
			bytes.NewReader([]byte(`{"name": "test", "latitude": 1.1, "longitude": 1.1}`)),
		)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		defer res.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func createMockService(t *testing.T) (*Mockactions, *gomock.Controller) {
	t.Helper()

	mockServiceController := gomock.NewController(t)
	mockService := NewMockactions(mockServiceController)

	return mockService, mockServiceController
}

func createTestApp() *fiber.App {
	return fiber.New(fiber.Config{})
}
