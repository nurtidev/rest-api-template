package main

import (
	"flag"
	"github.com/nurtidev/rest-api-template/internal/config"
	"github.com/nurtidev/rest-api-template/internal/server"
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

	srv, err := server.New(cfg, logger)
	if err != nil {
		log.Fatalf("init server: %v\n", err)
	}

	if err = srv.StartWithGracefulShutdown(); err != nil {
		log.Fatalf("start server: %v\n", err)
	}
}
