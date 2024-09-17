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
	GetUser(username string) (*domain.User, error)
}
