package user

import (
	"errors"

	"github.com/kaa-it/gophkeeper/internal/server/domain"
)

var (
	// ErrUserAlreadyExists indicates that a user with the same login already exists in the system.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrUserNotFound indicates that a user was not found in the system.
	ErrUserNotFound = errors.New("user not found")
)

// Repository provides methods for saving and retrieving User entities.
type Repository interface {
	// Save persists a User entity into the repository. Returns an error if the operation fails.
	Save(user *domain.User) error
	// Get retrieves a User entity based on the provided login. Returns the User if found, otherwise returns an error.
	Get(login string) (*domain.User, error)
}
