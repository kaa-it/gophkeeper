package credentials

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/kaa-it/gophkeeper/internal/server/domain"
)

// ErrCredentialsNotFound indicates that the requested credentials
// were not found in the storage.
var ErrCredentialsNotFound = errors.New("credentials not found")

// InMemoryCredentialsStore provides a thread-safe in-memory storage for credentials.
type InMemoryCredentialsStore struct {
	mutex       sync.RWMutex
	credentials map[string]*domain.Credentials
}

// NewInMemoryCredentialsStore initializes and returns a new instance
// of InMemoryCredentialsStore with empty credentials.
func NewInMemoryCredentialsStore() *InMemoryCredentialsStore {
	return &InMemoryCredentialsStore{
		credentials: make(map[string]*domain.Credentials),
	}
}

// Save stores the given credentials in the in-memory store and returns a generated unique identifier or an error.
func (s *InMemoryCredentialsStore) Save(credentials *domain.Credentials) (string, error) {
	credentialsID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate credentials id: %w", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.credentials[credentialsID.String()] = credentials.Clone()

	return credentialsID.String(), nil
}

// Get retrieves credentials by a given identifier from the in-memory store.
// Returns the corresponding credentials or an error if not found.
func (s *InMemoryCredentialsStore) Get(id string) (*domain.Credentials, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	credentials, ok := s.credentials[id]
	if !ok {
		return nil, ErrCredentialsNotFound
	}

	return credentials.Clone(), nil
}
