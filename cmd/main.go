package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nurtidev/rest-api-template/internal/config"
	"github.com/nurtidev/rest-api-template/internal/handler"
	"github.com/nurtidev/rest-api-template/internal/repository/postgres"
	"github.com/nurtidev/rest-api-template/internal/service"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "init config")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // todo: прокинуть с config file
	}))

	dsn := fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=disable", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Database)
	db, err := postgres.New(dsn, logger)
	if err != nil {
		return errors.Wrap(err, "init db")
	}

	svc, err := service.New(db, logger)
	if err != nil {
		return errors.Wrap(err, "init service")
	}

	h, err := handler.New(svc, logger)
	if err != nil {
		return errors.Wrap(err, "init handler")
	}

	app := fiber.New(fiber.Config{AppName: cfg.App.Name})
	h.Routes(app)

	return app.Listen(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
}
