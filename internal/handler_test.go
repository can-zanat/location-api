package internal

import (
	"bytes"
	"location-api/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/gofiber/fiber/v2"
)

func TestHandler_CreateLocation(t *testing.T) {
	t.Run("should create location properly", func(t *testing.T) {
		mockService, mockServiceController := createMockService(t)
		defer mockServiceController.Finish()

		app := createTestApp()

		mockService.
			EXPECT().
			CreateLocation(&testCreateLocationReq).
			Return(&testCreateLocationRes, nil).
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
			CreateLocation(&testCreateLocationReq).
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

func TestHandler_GetLocation(t *testing.T) {
	t.Run("should get location properly", func(t *testing.T) {
		mockService, mockServiceController := createMockService(t)
		defer mockServiceController.Finish()

		app := createTestApp()

		mockService.
			EXPECT().
			GetLocation(&testGetLocationReq).
			Return(&testGetLocationRes, nil).
			Times(1)

		handler := NewHandler(mockService)
		handler.RegisterRoutes(app)

		req := httptest.NewRequest(
			http.MethodGet,
			"/location?id=67d562e3d9f2d225ca4d9918",
			http.NoBody,
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
			GetLocation(&testGetLocationReq).
			Return(nil, assert.AnError).
			Times(1)

		handler := NewHandler(mockService)
		handler.RegisterRoutes(app)

		req := httptest.NewRequest(
			http.MethodGet,
			"/location?id=67d562e3d9f2d225ca4d9918",
			http.NoBody,
		)

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
			http.MethodGet,
			"/location?id=test",
			http.NoBody,
		)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		defer res.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestHandler_GetLocations(t *testing.T) {
	t.Run("should get location properly", func(t *testing.T) {
		mockService, mockServiceController := createMockService(t)
		defer mockServiceController.Finish()

		app := createTestApp()

		mockService.
			EXPECT().
			GetLocations(&testGetLocationsReq).
			Return(&testGetLocationsRes, nil).
			Times(1)

		handler := NewHandler(mockService)
		handler.RegisterRoutes(app)

		req := httptest.NewRequest(
			http.MethodGet,
			"/locations?page=1&limit=1",
			http.NoBody,
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
			GetLocations(&testGetLocationsReq).
			Return(nil, assert.AnError).
			Times(1)

		handler := NewHandler(mockService)
		handler.RegisterRoutes(app)

		req := httptest.NewRequest(
			http.MethodGet,
			"/locations?page=1&limit=1",
			http.NoBody,
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
			http.MethodGet,
			"/locations?page=1&limit=test",
			http.NoBody,
		)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		defer res.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestHandler_UpdateLocations(t *testing.T) {
	t.Run("should update locations properly", func(t *testing.T) {
		mockService, mockServiceController := createMockService(t)
		defer mockServiceController.Finish()

		app := createTestApp()

		mockService.
			EXPECT().
			UpdateLocations(&testUpdateLocationsReq).
			Return(&testUpdateLocationsRes, nil).
			Times(1)

		handler := NewHandler(mockService)
		handler.RegisterRoutes(app)

		reqBody := `{
			"locations": [
				{
					"id": "67d562e3d9f2d225ca4d9918",
					"name": "test",
					"latitude": 1.1,
					"longitude": 1.1,
					"marker_color": "FFFFFF"
				},
				{
					"id": "67d562e3d9f2d225ca4d9919",
					"name": "test2",
					"latitude": 2.2,
					"longitude": 2.2,
					"marker_color": "000000"
				}
			]
		}`

		req := httptest.NewRequest(
			http.MethodPatch,
			"/locations",
			bytes.NewReader([]byte(reqBody)),
		)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		defer res.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("should return partial content", func(t *testing.T) {
		mockService, mockServiceController := createMockService(t)
		defer mockServiceController.Finish()

		app := createTestApp()

		mockService.
			EXPECT().
			UpdateLocations(&testUpdateLocationsReqForPartialContent).
			Return(&model.UpdateLocationsResponse{
				UpdatedIDs:   []string{"67d562e3d9f2d225ca4d9918", "67d562e3d9f2d225ca4d9919"},
				FailedIDs:    []string{"67d562e3d9ddd225ca4d9919"},
				UpdatedCount: 2,
			}, nil).
			Times(1)

		handler := NewHandler(mockService)
		handler.RegisterRoutes(app)

		reqBody := `{
			"locations": [
				{
					"id": "67d562e3d9f2d225ca4d9918",
					"name": "test",
					"latitude": 1.1,
					"longitude": 1.1,
					"marker_color": "FFFFFF"
				},
				{
					"id": "67d562e3d9f2d225ca4d9919",
					"name": "test2",
					"latitude": 2.2,
					"longitude": 2.2,
					"marker_color": "000000"
				},
				{
					"id": "67d562e3d9ddd225ca4d9919",
					"name": "test3",
					"latitude": 3.3,
					"longitude": 3.3,
					"marker_color": "000000"
				}
			]
		}`

		req := httptest.NewRequest(
			http.MethodPatch,
			"/locations",
			bytes.NewReader([]byte(reqBody)),
		)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		defer res.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusPartialContent, res.StatusCode)
	})

	t.Run("should return internal server error", func(t *testing.T) {
		mockService, mockServiceController := createMockService(t)
		defer mockServiceController.Finish()

		app := createTestApp()

		mockService.
			EXPECT().
			UpdateLocations(&testUpdateLocationsReq).
			Return(nil, assert.AnError).
			Times(1)

		handler := NewHandler(mockService)
		handler.RegisterRoutes(app)

		reqBody := `{
			"locations": [
				{
					"id": "67d562e3d9f2d225ca4d9918",
					"name": "test",
					"latitude": 1.1,
					"longitude": 1.1,
					"marker_color": "FFFFFF"
				},
				{
					"id": "67d562e3d9f2d225ca4d9919",
					"name": "test2",
					"latitude": 2.2,
					"longitude": 2.2,	
					"marker_color": "000000"
				}
			]
		}`

		req := httptest.NewRequest(
			http.MethodPatch,
			"/locations",
			bytes.NewReader([]byte(reqBody)),
		)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		defer res.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("should return bad request error when body is empty", func(t *testing.T) {
		mockService, mockServiceController := createMockService(t)
		defer mockServiceController.Finish()

		app := createTestApp()

		handler := NewHandler(mockService)
		handler.RegisterRoutes(app)

		req := httptest.NewRequest(
			http.MethodPatch,
			"/locations",
			http.NoBody,
		)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req)
		defer res.Body.Close()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}

func TestHandler_GetRoutes(t *testing.T) {
	t.Run("should get routes properly", func(t *testing.T) {
		mockService, mockServiceController := createMockService(t)
		defer mockServiceController.Finish()

		app := createTestApp()

		mockService.
			EXPECT().
			GetRoutes(&testGetRoutesReq).
			Return(&testGetRoutesRes, nil).
			Times(1)

		handler := NewHandler(mockService)
		handler.RegisterRoutes(app)

		req := httptest.NewRequest(
			http.MethodGet,
			"/routes?longitude=1.1&latitude=1.1",
			http.NoBody,
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
			GetRoutes(&testGetRoutesReq).
			Return(nil, assert.AnError).
			Times(1)

		handler := NewHandler(mockService)
		handler.RegisterRoutes(app)

		req := httptest.NewRequest(
			http.MethodGet,
			"/routes?longitude=1.1&latitude=1.1",
			http.NoBody,
		)

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
			http.MethodGet,
			"/routes?longitude=test&latitude=1.1",
			http.NoBody,
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
