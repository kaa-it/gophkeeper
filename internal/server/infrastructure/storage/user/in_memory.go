package user

import (
	"sync"

	"github.com/kaa-it/gophkeeper/internal/server/domain"
)

// InMemoryUserStore represents an in-memory storage for User entities.
type InMemoryUserStore struct {
	mutex sync.RWMutex
	users map[string]*domain.User
}

// NewInMemoryUserStore creates and returns a new instance of InMemoryUserStore with initialized user map.
func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*domain.User),
	}
}

// Save stores a new User in the InMemoryUserStore. Returns an error if a user with the same login already exists.
func (s *InMemoryUserStore) Save(user *domain.User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.users[user.Login] != nil {
		return ErrUserAlreadyExists
	}

	s.users[user.Login] = user.Clone()
	return nil
}

// Get retrieves a User from the InMemoryUserStore based on the provided login.
// Returns the User or an error if not found.
func (s *InMemoryUserStore) Get(login string) (*domain.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user := s.users[login]
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user.Clone(), nil
}
