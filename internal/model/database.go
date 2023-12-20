package model

import (
	"time"
)

type Token struct {
	ID        int       `db:"id"`         // идентификатор токена
	UserID    int       `db:"user_id"`    // идентификатор пользователя которому принадлежит этот токен
	Value     string    `db:"value"`      // значение токена
	ExpiredAt time.Time `db:"expired_at"` // срок жизни токена
	CreatedAt time.Time `db:"created_at"` // время создания записи в базе
}

type User struct {
	Id             int       `db:"id"`              // идентификатор пользователя
	Email          string    `db:"email"`           // почта пользователя
	Name           string    `db:"name"`            // имя пользователя
	Surname        string    `db:"surname"`         // фамилия пользователя
	HashedPassword string    `db:"hashed_password"` // хэш пароля пользователя
	IsActive       bool      `db:"is_active"`       // статус активности аккаунта пользователя
	CreatedAt      time.Time `db:"created_at"`      // время создания записи в базе
	UpdatedAt      time.Time `db:"updated_at"`      // время обновления записи в базе
}
