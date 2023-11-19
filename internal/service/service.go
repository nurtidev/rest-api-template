package service

import (
	"github.com/nurtidev/rest-api-template/internal/repository"
	"log/slog"
)

// todo: подумать над тем нужен ли здесь интерфейс?

type Service interface {
}

type service struct {
	repo   repository.Repository
	logger *slog.Logger
}

func New(repo repository.Repository, logger *slog.Logger) (Service, error) {
	return &service{repo: repo, logger: logger}, nil
}
