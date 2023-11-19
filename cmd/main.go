package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nurtidev/rest-api-template/internal/config"
	"github.com/nurtidev/rest-api-template/internal/handler"
	"github.com/nurtidev/rest-api-template/internal/repository/postgres"
	"github.com/nurtidev/rest-api-template/internal/service"
	"log"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	if err := run(logger); err != nil {
		log.Fatal(err)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Init("./")
	if err != nil {
		return err
	}

	db, err := postgres.New(fmt.Sprintf("%s%s", cfg.Postgres.Host, cfg.Postgres.Port))
	if err != nil {
		return err
	}

	svc, err := service.New(db, logger)
	if err != nil {
		return err
	}

	h, err := handler.New(svc, logger)
	if err != nil {
		return err
	}

	app := fiber.New(fiber.Config{AppName: "backend.api"})
	h.Routes(app)

	return nil
}
