package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nurtidev/rest-api-template/internal/model"
	"time"
)

const defaultTimeout = 3 * time.Second

type Repository struct {
	db *sqlx.DB
}

func (r *Repository) FindUser(ctx context.Context, id int) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) InsertUser(ctx context.Context, u *model.User) (int, error) {
	//TODO implement me
	panic("implement me")
}

func New(dsn string) (*Repository, error) {
	db, err := sqlx.Connect("postgres", "postgres://"+dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(2 * time.Hour)

	return &Repository{db: db}, nil
}
