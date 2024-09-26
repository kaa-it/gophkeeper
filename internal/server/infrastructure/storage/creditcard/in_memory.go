package creditcard

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/kaa-it/gophkeeper/internal/server/domain"
)

// ErrCreditCardNotFound is returned when a specified credit card cannot be found in the storage.
var ErrCreditCardNotFound = errors.New("credit card not found")

// InMemoryCreditCardStore is an in-memory storage implementation for credit card data.
type InMemoryCreditCardStore struct {
	mutex       sync.RWMutex
	creditCards map[string]*domain.CreditCard
}

// NewInMemoryCreditCardStore initializes and returns a new instance
// of InMemoryCreditCardStore with an empty credit card map.
func NewInMemoryCreditCardStore() *InMemoryCreditCardStore {
	return &InMemoryCreditCardStore{
		creditCards: make(map[string]*domain.CreditCard),
	}
}

// Save stores a new credit card into the in-memory credit card store and returns the generated credit card ID.
func (s *InMemoryCreditCardStore) Save(creditCard *domain.CreditCard) (string, error) {
	creditCardID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("could not generate new creditCard id: %w", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.creditCards[creditCardID.String()] = creditCard

	return creditCardID.String(), nil
}
