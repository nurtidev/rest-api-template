package postgres

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/nurtidev/rest-api-template/internal/model"
)

func (r *Repository) FindToken(ctx context.Context, token *model.Token) (*model.Token, error) {
	queryBuilder := sq.Select("*").From(defaultSchema + tokensTable).PlaceholderFormat(sq.Dollar)

	if token.UserID != 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"user_id": token.UserID})
	}

	if token.Value != "" {
		queryBuilder = queryBuilder.Where(sq.Eq{"value": token.Value})
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var foundToken model.Token
	err = r.db.GetContext(ctx, &foundToken, query, args...)
	if err != nil {
		return nil, err
	}

	return &foundToken, nil
}

func (r *Repository) InsertToken(ctx context.Context, token *model.Token) (int, error) {
	queryBuilder := sq.Insert(defaultSchema+tokensTable).
		Columns("user_id", "value", "expired_at", "created_at").
		Values(token.UserID, token.Value, token.ExpiredAt, token.CreatedAt).
		Suffix("RETURNING id").PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var tokenID int
	err = r.db.GetContext(ctx, &tokenID, query, args...)
	if err != nil {
		return 0, err
	}

	return tokenID, nil
}

func (r *Repository) DeleteToken(ctx context.Context, token *model.Token) error {
	queryBuilder := sq.Delete(defaultSchema + tokensTable).PlaceholderFormat(sq.Dollar)
	if token.ID != 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"id": token.ID})
	}
	if token.Value != "" {
		queryBuilder = queryBuilder.Where(sq.Eq{"value": token.Value})
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
