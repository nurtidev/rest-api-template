package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nurtidev/rest-api-template/internal/model"
	"log/slog"
	"time"
)

const (
	defaultTimeout  = 3 * time.Second
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxIdleTime = 5 * time.Minute
	connMaxLifetime = 2 * time.Hour
)

type Repository struct {
	logger *slog.Logger
	db     *sqlx.DB
}

func (r *Repository) FindUser(ctx context.Context, id int) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) InsertUser(ctx context.Context, u *model.User) (int, error) {
	//TODO implement me
	panic("implement me")
}

func New(dsn string, logger *slog.Logger) (*Repository, error) {
	db, err := sqlx.Connect("postgres", "postgres://"+dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	db.SetConnMaxLifetime(connMaxLifetime)

	return &Repository{db: db, logger: logger}, nil
}
