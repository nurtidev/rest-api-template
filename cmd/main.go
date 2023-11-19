package main

import (
	"flag"
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
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	confPath := flag.String("config file path", "./configs/", "Path to configuration file")
	flag.Parse()

	cfg, err := config.Init(*confPath)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // todo: прокинуть с config file
	}))

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
