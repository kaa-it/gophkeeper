package user

import (
	"errors"

	"github.com/kaa-it/gophkeeper/internal/server/domain"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type Repository interface {
	Save(user *domain.User) error
	GetUser(login string) (*domain.User, error)
}
