package creditcard

import "github.com/kaa-it/gophkeeper/internal/server/domain"

// Repository is an interface for saving CreditCard information.
type Repository interface {
	// Save stores the given credit card information in the repository and returns a unique identifier or an error.
	Save(creditCard *domain.CreditCard) (string, error)
}
