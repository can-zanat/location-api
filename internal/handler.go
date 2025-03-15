package internal

import (
	"location-api/model"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service actions
}

type actions interface {
	CreateLocation(req *model.CreateLocationRequest) (*model.CreateLocationResponse, error)
}

func NewHandler(service actions) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Post("/location", h.CreateLocation)
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
