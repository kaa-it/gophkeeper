package user

import (
	"sync"

	"github.com/kaa-it/gophkeeper/internal/server/domain"
)

type InMemoryUserStore struct {
	mutex sync.RWMutex
	users map[string]*domain.User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*domain.User),
	}
}

func (s *InMemoryUserStore) Save(user *domain.User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.users[user.Login] != nil {
		return ErrUserAlreadyExists
	}

	s.users[user.Login] = user.Clone()
	return nil
}

func (s *InMemoryUserStore) GetUser(login string) (*domain.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user := s.users[login]
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user.Clone(), nil
}
