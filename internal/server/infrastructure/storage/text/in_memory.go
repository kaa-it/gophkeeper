package text

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/kaa-it/gophkeeper/internal/server/domain"
)

// ErrTextNotFound represents an error indicating that the specified text was not found.
var ErrTextNotFound = errors.New("text not found")

// InMemoryTextStore is a thread-safe in-memory storage for text objects.
type InMemoryTextStore struct {
	mutex sync.RWMutex
	texts map[string]*domain.Text
}

// NewInMemoryTextStore initializes and returns a new instance of InMemoryTextStorage with an empty text map.
func NewInMemoryTextStore() *InMemoryTextStore {
	return &InMemoryTextStore{
		texts: make(map[string]*domain.Text),
	}
}

// Save stores a Text object in the in-memory storage and returns a unique text ID
// or an error if the ID generation fails.
func (s *InMemoryTextStore) Save(text *domain.Text) (string, error) {
	textID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("unable to generate textID: %w", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.texts[textID.String()] = text

	return textID.String(), nil
}
