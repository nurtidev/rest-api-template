package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nurtidev/rest-api-template/internal/service"
	"log/slog"
)

type Handler struct {
	service service.Service
	logger  *slog.Logger
}

func New(s service.Service, logger *slog.Logger) (Handler, error) {
	return Handler{service: s, logger: logger}, nil
}

func (h *Handler) Routes(app *fiber.App) {
	router := app.Group("/api/v1")

	router.Get("/ping", h.Ping)
}

func (h *Handler) Ping(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": "pong",
	})
}
