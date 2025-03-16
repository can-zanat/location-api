package internal

import (
	"location-api/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service actions
}

type actions interface {
	CreateLocation(req *model.CreateLocationRequest) (*model.CreateLocationResponse, error)
	GetLocation(req *model.GetLocationRequest) (*model.GetLocationResponse, error)
	GetLocations(req *model.GetLocationsRequest) (*model.GetLocationsResponse, error)
	UpdateLocations(req *model.UpdateLocationsRequest) (*model.UpdateLocationsResponse, error)
	GetRoutes(req *model.GetRoutesRequest) (*model.GetRoutesResponse, error)
}

func NewHandler(service actions) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Post("/location", h.CreateLocation)
	app.Get("/location", h.GetLocation)
	app.Get("/locations", h.GetLocations)
	app.Patch("/locations", h.UpdateLocations)
	app.Get("/routes", h.GetRoutes)
}

func (h *Handler) CreateLocation(ctx *fiber.Ctx) error {
	var req model.CreateLocationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := req.ValidateLocation(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.service.CreateLocation(&req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *Handler) GetLocation(ctx *fiber.Ctx) error {
	var req model.GetLocationRequest

	if err := ctx.QueryParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	res, err := h.service.GetLocation(&req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *Handler) GetLocations(ctx *fiber.Ctx) error {
	var req model.GetLocationsRequest

	if err := ctx.QueryParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	res, err := h.service.GetLocations(&req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *Handler) UpdateLocations(ctx *fiber.Ctx) error {
	var req model.UpdateLocationsRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if len(req.Locations) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No locations provided",
		})
	}

	if err := req.ValidateLocation(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.service.UpdateLocations(&req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if len(res.FailedIDs) > 0 && len(res.UpdatedIDs) > 0 {
		return ctx.Status(fiber.StatusPartialContent).JSON(res)
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *Handler) GetRoutes(ctx *fiber.Ctx) error {
	var req model.GetRoutesRequest

	if err := ctx.QueryParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := req.ValidateLocation(); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.service.GetRoutes(&req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
