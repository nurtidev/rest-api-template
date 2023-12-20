package service

import (
	"context"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/nurtidev/rest-api-template/internal/config"
	"github.com/nurtidev/rest-api-template/internal/model"
	"github.com/nurtidev/rest-api-template/internal/repository"
	"github.com/pkg/errors"
	"log/slog"
	"time"
)

// todo: подумать над тем нужен ли здесь интерфейс?

type Service interface {
	RegisterUser(ctx context.Context, input *model.RegisterRequest) (string, int, error)
	LoginUser(ctx context.Context, input *model.LoginRequest) (string, int, error)
	LogoutUser(ctx context.Context, token string) (int, error)
	RefreshToken(ctx context.Context, token string) (string, int, error)
	ValidateToken(ctx context.Context, token string) (int, error)
}

type service struct {
	cfg    *config.Config
	repo   repository.Repository
	logger *slog.Logger
}

func New(cfg *config.Config, repo repository.Repository, logger *slog.Logger) (Service, error) {
	return &service{cfg: cfg, repo: repo, logger: logger}, nil
}

func (s *service) ValidateToken(ctx context.Context, token string) (int, error) {
	t, err := validateJWTToken([]byte(s.cfg.Secrets.JwtSecret), token)
	if err != nil {
		return fiber.StatusBadRequest, errors.Wrap(err, "invalid token")
	}

	_, err = s.repo.FindToken(ctx, &model.Token{
		UserID: t.UserID,
		Value:  token,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fiber.StatusInternalServerError, errors.Wrap(err, "delete token by value")
	} else if errors.Is(err, sql.ErrNoRows) {
		return fiber.StatusBadRequest, errors.Wrap(err, "token not found")
	}

	return fiber.StatusOK, nil
}

func (s *service) RefreshToken(ctx context.Context, token string) (string, int, error) {
	t, err := validateJWTToken([]byte(s.cfg.Secrets.JwtSecret), token)
	if err != nil {
		return "", fiber.StatusUnauthorized, errors.Wrap(err, "invalid token")
	}

	_, err = s.repo.FindToken(ctx, &model.Token{Value: token})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fiber.StatusInternalServerError, errors.Wrap(err, "find token by value")
	} else if errors.Is(err, sql.ErrNoRows) {
		return "", fiber.StatusBadRequest, errors.Wrap(err, "token not found")
	}

	err = s.repo.DeleteToken(ctx, &model.Token{Value: token})
	if err != nil {
		return "", fiber.StatusInternalServerError, errors.Wrap(err, "delete token by value")
	}

	token, expiredAt, err := generateJWTToken(t.UserID, t.Email, s.cfg.Secrets.JwtSecret)
	if err != nil {
		return "", fiber.StatusInternalServerError, errors.Wrap(err, "generate token")
	}

	_, err = s.repo.InsertToken(ctx, &model.Token{
		UserID:    t.UserID,
		Value:     token,
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return "", fiber.StatusInternalServerError, errors.Wrap(err, "insert token")
	}

	return token, fiber.StatusOK, nil
}

func (s *service) LogoutUser(ctx context.Context, token string) (int, error) {
	u, err := validateJWTToken([]byte(s.cfg.Secrets.JwtSecret), token)
	if err != nil {
		return fiber.StatusBadRequest, errors.Wrap(err, "invalid token")
	}

	err = s.repo.DeleteToken(ctx, &model.Token{Value: token})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fiber.StatusInternalServerError, errors.Wrap(err, "delete token by value")
	} else if errors.Is(err, sql.ErrNoRows) {
		return fiber.StatusBadRequest, errors.Wrap(err, "token not found")
	}
	s.logger.Info("user successfully logout", u.Email)
	return fiber.StatusOK, nil
}

func (s *service) LoginUser(ctx context.Context, input *model.LoginRequest) (string, int, error) {
	u, err := s.repo.FindUser(ctx, &model.User{
		Email: input.Email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fiber.StatusUnauthorized, errors.New("invalid login or password")
		}
		s.logger.Error("find user by email", input.Email, err)
		return "", fiber.StatusInternalServerError, err
	}

	if err = comparePassword(input.Password, u.HashedPassword); err != nil {
		s.logger.Error("compare password", u.Email, err)
		return "", fiber.StatusUnauthorized, errors.New("invalid login or password")
	}

	tokenModel := &model.Token{}
	tokenModel, err = s.repo.FindToken(ctx, &model.Token{UserID: u.Id})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Error("find token by user id", u.Email, err)
		return "", fiber.StatusInternalServerError, err
	} else if errors.Is(err, sql.ErrNoRows) {
		token, expiredAt, err := generateJWTToken(u.Id, u.Email, s.cfg.Secrets.JwtSecret)
		if err != nil {
			s.logger.Error("generate jwt new token by user", u.Email, err)
			return "", fiber.StatusInternalServerError, errors.Wrap(err, "generate token")
		}

		_, err = s.repo.InsertToken(ctx, &model.Token{
			UserID:    u.Id,
			Value:     token,
			ExpiredAt: expiredAt,
			CreatedAt: time.Now(),
		})
		if err != nil {
			s.logger.Error("insert token into db", u.Email, err)
			return "", fiber.StatusInternalServerError, errors.Wrap(err, "insert token")
		}

		tokenModel.Value = token
	}

	s.logger.Info("user successfully login", u.Email)
	return tokenModel.Value, fiber.StatusOK, nil
}

func (s *service) RegisterUser(ctx context.Context, input *model.RegisterRequest) (string, int, error) {
	u, err := s.repo.FindUser(ctx, &model.User{
		Email: input.Email,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Error("find user by param", err)
		return "", fiber.StatusInternalServerError, errors.Wrap(err, "find user by param")
	}

	if u != nil {
		s.logger.Error("duplicate user email", input.Email)
		return "", fiber.StatusBadRequest, errors.New("duplicate user email")
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		s.logger.Error("hash password", input.Email, input.Password, err)
		return "", fiber.StatusInternalServerError, errors.Wrap(err, "hash password")
	}

	id, err := s.repo.InsertUser(ctx, &model.User{
		Email:          input.Email,
		Name:           input.Name,
		Surname:        input.Surname,
		HashedPassword: hashedPassword,
		IsActive:       true, // TODO: потом отдельно вынести активацию пользователя
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		s.logger.Error("insert user into db", input.Email, err)
		return "", fiber.StatusInternalServerError, errors.Wrap(err, "insert user")
	}

	s.logger.Info("user successfully created", input.Email)

	token, expiredAt, err := generateJWTToken(id, input.Email, s.cfg.Secrets.JwtSecret)
	if err != nil {
		s.logger.Error("generate token by user", input.Email, err)
		return "", fiber.StatusInternalServerError, errors.Wrap(err, "generate token")
	}

	_, err = s.repo.InsertToken(ctx, &model.Token{
		UserID:    id,
		Value:     token,
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
	})
	if err != nil {
		s.logger.Error("insert token into db", input.Email, err)
		return "", fiber.StatusInternalServerError, errors.Wrap(err, "insert token")
	}

	return token, fiber.StatusOK, nil

}
