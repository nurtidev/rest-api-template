package repository

import (
	"context"
	"github.com/nurtidev/rest-api-template/internal/model"
)

type Repository interface {
	FindUser(ctx context.Context, user *model.User) (*model.User, error)
	InsertUser(ctx context.Context, user *model.User) (int, error)
	UpdateUser(ctx context.Context, user *model.User) error

	FindToken(ctx context.Context, token *model.Token) (*model.Token, error)
	InsertToken(ctx context.Context, token *model.Token) (int, error)
	DeleteToken(ctx context.Context, token *model.Token) error

	Close(ctx context.Context) error
}
