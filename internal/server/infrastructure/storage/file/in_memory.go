package file

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

type InMemoryFileStore struct {
	mutex         sync.Mutex
	storageFolder string
	files         map[string]*Info
}

type Info struct {
	Metadata string
	FileName string
	Path     string
}

func NewInMemoryFileStore(storageFolder string) *InMemoryFileStore {
	return &InMemoryFileStore{
		storageFolder: storageFolder,
		files:         make(map[string]*Info),
	}
}

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
