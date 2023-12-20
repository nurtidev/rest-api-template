package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nurtidev/rest-api-template/internal/model"
	"strings"
)

func (h *Handler) recoverPanic(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Info("recovered from panic", r)
			_ = c.Status(fiber.StatusInternalServerError).SendString("internal server error")
		}
	}()

	return c.Next()
}

func (h *Handler) protected(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).SendString("missing or malformed auth token")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	st, err := h.service.ValidateToken(c.UserContext(), token)
	if err != nil {
		return c.Status(st).JSON(model.BaseResponse{
			Success: false,
			Msg:     err.Error(),
		})
	}

	return c.Next()
}
