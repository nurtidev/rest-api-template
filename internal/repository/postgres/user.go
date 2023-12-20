package postgres

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/nurtidev/rest-api-template/internal/model"
)

func (r *Repository) FindUser(ctx context.Context, user *model.User) (*model.User, error) {
	users := sq.Select("*").From(defaultSchema + usersTable).PlaceholderFormat(sq.Dollar)

	if user.Id != 0 {
		users = users.Where(sq.Eq{"id": user.Id})
	}

	if user.Email != "" {
		users = users.Where(sq.Eq{"email": user.Email})
	}

	query, args, err := users.ToSql()
	if err != nil {
		return nil, err
	}

	var foundUser model.User
	err = r.db.GetContext(ctx, &foundUser, query, args...)
	if err != nil {
		return nil, err
	}

	return &foundUser, nil
}

func (r *Repository) InsertUser(ctx context.Context, user *model.User) (int, error) {
	queryBuilder := sq.Insert(defaultSchema+usersTable).
		Columns("email", "name", "surname", "hashed_password", "is_active", "created_at", "updated_at").
		Values(user.Email, user.Name, user.Surname, user.HashedPassword, user.IsActive, user.CreatedAt, user.UpdatedAt).
		Suffix("RETURNING id").PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var userID int
	err = r.db.GetContext(ctx, &userID, query, args...)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *model.User) error {
	queryBuilder := sq.Update(defaultSchema+usersTable).
		Set("email", user.Email).
		Set("name", user.Name).
		Set("hashed_password", user.HashedPassword).
		Set("is_active", user.IsActive).
		Set("updated_at", user.UpdatedAt).
		Where(sq.Eq{"id": user.Id}).PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
