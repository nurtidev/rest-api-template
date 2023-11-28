package main

import (
	"flag"
	"fmt"
	"github.com/nurtidev/rest-api-template/internal/config"
	"github.com/nurtidev/rest-api-template/internal/handler"
	"github.com/nurtidev/rest-api-template/internal/repository/mock"
	"github.com/nurtidev/rest-api-template/internal/service"
	"github.com/nurtidev/rest-api-template/internal/utils"
	"log"
	"log/slog"
	"os"
)

var confPath = flag.String("config file path", "./configs/", "Path to configuration file")

func main() {
	flag.Parse()

	cfg, err := config.Init(*confPath)
	if err != nil {
		log.Fatalf("init config: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: utils.ToSlogLevel(cfg.Logger.Level),
	}))

	dsn := fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=disable", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Database)
	db, err := mock.New(dsn, logger)
	if err != nil {
		log.Fatalf("init db: %v", err)
	}

	svc, err := service.New(db, logger)
	if err != nil {
		log.Fatalf("init service: %v", err)
	}

	h, err := handler.New(svc, logger)
	if err != nil {
		log.Fatalf("init handler: %v", err)
	}

	if err = h.ServeHTTP(cfg); err != nil {
		log.Fatalf("serve http: %v", err)
	}
}
