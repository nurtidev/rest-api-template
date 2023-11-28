package mock

import (
	"context"
	"github.com/nurtidev/rest-api-template/internal/model"
	"log/slog"
)

type Repository struct {
	logger *slog.Logger
}

func (r *Repository) FindUser(ctx context.Context, id int) (*model.User, error) {
	return &model.User{
		Login:   "mock_user",
		Address: "mock_address",
	}, nil
}

func (r *Repository) InsertUser(ctx context.Context, u *model.User) (int, error) {
	return 1, nil
}

func New(dsn string, logger *slog.Logger) (*Repository, error) {
	return &Repository{logger: logger}, nil
}
