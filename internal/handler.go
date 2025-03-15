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
}

func NewHandler(service actions) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Post("/location", h.CreateLocation)
	app.Get("/location", h.GetLocation)
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
