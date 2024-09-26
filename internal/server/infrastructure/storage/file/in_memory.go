package file

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

// InMemoryFileStore is a struct that provides an in-memory storage solution for file info with thread-safe operations.
type InMemoryFileStore struct {
	mutex         sync.Mutex
	storageFolder string
	files         map[string]*Info
}

// Info represents file information including metadata, file name, and path.
type Info struct {
	Metadata string
	FileName string
	Path     string
}

// NewInMemoryFileStore initializes a new InMemoryFileStore with the provided storage folder.
func NewInMemoryFileStore(storageFolder string) *InMemoryFileStore {
	return &InMemoryFileStore{
		storageFolder: storageFolder,
		files:         make(map[string]*Info),
	}
}

// Save stores the file data with the associated metadata and filename, returning a unique file ID or an error.
func (s *InMemoryFileStore) Save(metadata string, fileName string, fileData bytes.Buffer) (string, error) {
	fileID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate file id: %w", err)
	}

	filePath := fmt.Sprintf("%s/%s", s.storageFolder, fileID)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	_, err = fileData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("failed to write file data: %w", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.files[fileID.String()] = &Info{
		Metadata: metadata,
		FileName: fileName,
		Path:     filePath,
	}

	return fileID.String(), nil
}
