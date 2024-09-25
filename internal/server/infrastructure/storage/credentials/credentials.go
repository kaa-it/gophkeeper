package credentials

import "github.com/kaa-it/gophkeeper/internal/server/domain"

// Repository defines an interface for managing credential entities in storage.
type Repository interface {
	// Save method stores given credentials and returns a unique ID or an error.
	Save(credentials *domain.Credentials) (string, error)
	// Get method retrieves credentials by unique ID, returning the credentials or an error.
	Get(id string) (*domain.Credentials, error)
}
