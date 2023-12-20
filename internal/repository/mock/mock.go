package mock

import (
	"log/slog"
)

const (
	usersTable  = "users"
	tokensTable = "tokens"
)

type Repository struct {
	logger *slog.Logger
	db     map[string]map[int]interface{}
}

func New(dsn string, logger *slog.Logger) (*Repository, error) {
	db := make(map[string]map[int]interface{})
	return &Repository{logger: logger, db: db}, nil
}
