package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nurtidev/rest-api-template/internal/config"
	"log/slog"
	"time"
)

const (
	defaultTimeout  = 3 * time.Second
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxIdleTime = 5 * time.Minute
	connMaxLifetime = 2 * time.Hour

	defaultSchema = "test_nurtilek."

	usersTable  = "users"
	tokensTable = "tokens"
)

type Repository struct {
	logger *slog.Logger
	db     *sqlx.DB
}

func New(cfg *config.Config, logger *slog.Logger) (*Repository, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Database, cfg.Postgres.SslMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	db.SetConnMaxLifetime(connMaxLifetime)

	return &Repository{db: db, logger: logger}, nil
}

func (r *Repository) Close(ctx context.Context) error {
	dbCloseChan := make(chan error, 1)
	go func() {
		dbCloseChan <- r.db.Close()
	}()

	select {
	case err := <-dbCloseChan:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}
