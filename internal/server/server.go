package server

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nurtidev/rest-api-template/internal/config"
	"github.com/nurtidev/rest-api-template/internal/handler"
	"github.com/nurtidev/rest-api-template/internal/repository"
	"github.com/nurtidev/rest-api-template/internal/repository/postgres"
	"github.com/nurtidev/rest-api-template/internal/service"
	"github.com/pkg/errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

type Server struct {
	app    *fiber.App
	repo   repository.Repository
	cfg    *config.Config
	logger *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) (*Server, error) {
	app := fiber.New(fiber.Config{
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
		AppName:      cfg.App.Name,
	})

	repo, err := postgres.New(cfg, logger)
	if err != nil {
		return nil, errors.Wrap(err, "init repo")
	}

	svc, err := service.New(cfg, repo, logger)
	if err != nil {
		return nil, errors.Wrap(err, "init service")
	}

	h, err := handler.New(svc, logger)
	if err != nil {
		return nil, errors.Wrap(err, "init handler")
	}

	h.InitializeRoutes(app)

	return &Server{
		app:    app,
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (s *Server) Start() error {
	s.logger.Info("starting server", slog.Group("server", "addr", s.cfg.Server.Host+":"+s.cfg.Server.Port))

	err := s.app.Listen(fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port))
	if err != nil {
		return err
	}

	s.logger.Info("stopped server", slog.Group("server", "addr", s.cfg.Server.Host+":"+s.cfg.Server.Port))
	return nil
}

func (s *Server) StartWithGracefulShutdown() error {
	shutdownErrChan := make(chan error)

	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		if err := s.repo.Close(ctx); err != nil {
			s.logger.Error("error closing database", err)
			shutdownErrChan <- err
			return
		}

		if err := s.app.ShutdownWithContext(ctx); err != nil {
			s.logger.Error("error shutting down server", err)
			shutdownErrChan <- err
			return
		}

		close(shutdownErrChan)
	}()

	s.logger.Info("starting server", slog.Group("server", "addr", s.cfg.Server.Host+":"+s.cfg.Server.Port))

	err := s.app.Listen(fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port))
	if err != nil {
		return err
	}

	err = <-shutdownErrChan
	if err != nil {
		return err
	}

	s.logger.Info("stopped server", slog.Group("server", "addr", s.cfg.Server.Host+":"+s.cfg.Server.Port))
	return nil
}
