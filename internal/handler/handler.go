package handler

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/nurtidev/rest-api-template/internal/model"
	"github.com/nurtidev/rest-api-template/internal/service"
	"log/slog"
	"strings"
)

type Handler struct {
	service service.Service
	logger  *slog.Logger
}

func New(s service.Service, logger *slog.Logger) (Handler, error) {
	return Handler{service: s, logger: logger}, nil
}

// @title Swagger API
// @version 2.0
// @description This is a sample rest api template.

func (h *Handler) InitializeRoutes(app *fiber.App) {

	app.Use(h.recoverPanic)

	app.Get("health", h.health)
	app.Get("/swagger/*", swagger.HandlerDefault)

	router := app.Group("/api/v1")

	auth := router.Group("/auth")
	auth.Post("/login", h.login)
	auth.Post("/register", h.register)
	auth.Post("/refresh", h.refresh)
	auth.Post("/logout", h.logout)
}

// health show the status of server
// @Summary show the status of server.
// @Description get the status of server.
// @Produce json
// @Success 200 {object} model.HealthResponse
// @Router /health [get]
func (h *Handler) health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(model.HealthResponse{Status: "ok"})
}

// logout deletes the JWT token from storage.
// @Summary delete jwt token from storage.
// @Description delete jwt token from storage.
// @Produce json
// @Param Authorization header string true "Bearer [JWT token]"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 401 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/v1/auth/logout [post]
func (h *Handler) logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).SendString("missing or malformed auth token")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	st, err := h.service.LogoutUser(c.UserContext(), token)
	if err != nil {
		return c.Status(st).JSON(model.BaseResponse{
			Success: false,
			Msg:     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.BaseResponse{
		Success: true,
		Msg:     "user successfully logout",
	})
}

// login authenticates a user and returns a JWT token.
// @Summary User login
// @Description Authenticates user credentials and returns a JWT token upon successful login.
// @Accept json
// @Produce json
// @Param login body model.LoginRequest true "Login Credentials"
// @Success 200 {object} model.LoginResponse "User successfully logged in with token returned"
// @Failure 400 {object} model.BaseResponse "Invalid request format or content"
// @Failure 401 {object} model.BaseResponse "Unauthorized access due to invalid credentials"
// @Failure 500 {object} model.BaseResponse "Internal server error"
// @Router /api/v1/auth/login [post]
func (h *Handler) login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"msg":     err.Error(),
		})
	}

	token, st, err := h.service.LoginUser(c.UserContext(), &req)
	if err != nil {
		return c.Status(st).JSON(model.BaseResponse{
			Success: false,
			Msg:     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.LoginResponse{
		BaseResponse: model.BaseResponse{
			Success: true,
			Msg:     "user successfully login",
		},
		Token: token,
	})
}

// register creates a new user account and returns a JWT token.
// @Summary User registration
// @Description Registers a new user with the provided details and returns a JWT token upon successful registration.
// @Accept json
// @Produce json
// @Param register body model.RegisterRequest true "Registration Details"
// @Success 200 {object} model.RegisterResponse "User successfully created with token returned"
// @Failure 400 {object} model.BaseResponse "Invalid request format or content"
// @Failure 500 {object} model.BaseResponse "Internal server error"
// @Router /api/v1/auth/register [post]
func (h *Handler) register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.BaseResponse{
			Success: false,
			Msg:     err.Error(),
		})
	}
	token, st, err := h.service.RegisterUser(c.UserContext(), &req)
	if err != nil {
		return c.Status(st).JSON(model.BaseResponse{
			Success: false,
			Msg:     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.RegisterResponse{
		BaseResponse: model.BaseResponse{
			Success: true,
			Msg:     "user successfully created",
		},
		Token: token,
	})
}

// refresh renews the JWT token for a user.
// @Summary Refresh JWT token
// @Description Validates the existing JWT token from the Authorization header and issues a new token.
// @Produce json
// @Param Authorization header string true "Bearer [current JWT token]"
// @Success 200 {object} model.RefreshResponse "Token successfully refreshed"
// @Failure 400 {object} model.BaseResponse "Missing or malformed auth token"
// @Failure 401 {object} model.BaseResponse "Unauthorized access due to invalid or expired token"
// @Failure 500 {object} model.BaseResponse "Internal server error"
// @Router /api/v1/auth/refresh [post]
func (h *Handler) refresh(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(model.BaseResponse{
			Success: false,
			Msg:     "missing or malformed auth token",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	token, st, err := h.service.RefreshToken(c.UserContext(), token)
	if err != nil {
		return c.Status(st).JSON(model.BaseResponse{
			Success: false,
			Msg:     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.RefreshResponse{
		BaseResponse: model.BaseResponse{
			Success: true,
			Msg:     "token successfully refreshed",
		},
		Token: token,
	})
}
