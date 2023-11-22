package handler

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (h *Handler) recoverPanic(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Info("recovered from panic", r)
			_ = c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
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

	// todo: validate token
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid or expired auth token")
	}

	return c.Next()
}
