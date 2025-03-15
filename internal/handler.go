package internal

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service actions
}

type actions interface {
}

func NewHandler(service actions) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {}
