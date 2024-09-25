package text

import "github.com/kaa-it/gophkeeper/internal/server/domain"

// Repository defines an interface for saving domain.Text entities.
type Repository interface {
	// Save persists a given domain.Text entity to the repository.
	// It returns the ID of the saved entity or an error if the save operation fails.
	Save(text *domain.Text) (string, error)
}
