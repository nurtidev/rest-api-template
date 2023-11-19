package repository

import (
	"context"
	"github.com/nurtidev/rest-api-template/internal/model"
)

type Repository interface {
	FindUser(ctx context.Context, id int) (*model.User, error)
	InsertUser(ctx context.Context, u *model.User) (int, error)
}
