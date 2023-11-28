package handler

import (
	"context"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	_ "github.com/nurtidev/rest-api-template/docs"
	"github.com/nurtidev/rest-api-template/internal/config"
	"github.com/nurtidev/rest-api-template/internal/model"
	"github.com/nurtidev/rest-api-template/internal/service"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

type Handler struct {
	service service.Service
	logger  *slog.Logger
	wg      sync.WaitGroup
}

func New(s service.Service, logger *slog.Logger) (Handler, error) {
	return Handler{service: s, logger: logger, wg: sync.WaitGroup{}}, nil
}

// @title Swagger API
// @version 2.0
// @description This is a sample rest api template.
func (h *Handler) ServeHTTP(cfg *config.Config) error {
	app := fiber.New(fiber.Config{
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
		AppName:      cfg.App.Name,
	})

	shutdownErrChan := make(chan error)

	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		shutdownErrChan <- app.ShutdownWithContext(ctx)
	}()

	h.logger.Info("starting server", slog.Group("server", "addr", cfg.Server.Host+":"+cfg.Server.Port))

	h.routes(app)

	err := app.Listen(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		return err
	}

	err = <-shutdownErrChan
	if err != nil {
		return err
	}

	h.logger.Info("stopped server", slog.Group("server", "addr", cfg.Server.Host+":"+cfg.Server.Port))

	h.wg.Wait()
	return nil
}

func (h *Handler) routes(app *fiber.App) {

	app.Use(h.recoverPanic)

	app.Get("health", h.health)
	app.Get("/swagger/*", swagger.HandlerDefault)

	router := app.Group("/api/v1")

	auth := router.Group("/auth")
	auth.Post("/login", h.login)
	auth.Post("/register", h.register)
	auth.Post("/refresh", h.refresh)

	protected := router.Group("/protected").Use(h.protected)
	protected.Get("/health", h.health)
}

// Health godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Produce json
// @Success 200 {object} model.HealthResponse
// @Router /health [get]
func (h *Handler) health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(model.HealthResponse{Status: "ok"})
}

func (h *Handler) login(c *fiber.Ctx) error {
	type request struct {
		email    string
		password string
	}
	var req request
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	type token struct {
		value     string
		expiredAt time.Time
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": "user successfully login",
		"token": &token{
			value:     "token",
			expiredAt: time.Now().Add(1 * time.Hour),
		},
	})
}

func (h *Handler) register(c *fiber.Ctx) error {
	type request struct {
		email    string
		password string
	}
	var req request
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	type token struct {
		value     string
		expiredAt time.Time
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": "user successfully created",
		"token": &token{
			value:     "token",
			expiredAt: time.Now().Add(1 * time.Hour),
		},
	})
}

func (h *Handler) refresh(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).SendString("missing or malformed auth token")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	// todo: validate token
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid or expired auth token")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg":   "token successfully refreshed",
		"token": &token,
	})
}
