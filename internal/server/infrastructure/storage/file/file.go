package file

import "bytes"

// Repository defines an interface for saving file metadata and data into storage.
type Repository interface {
	// Save method stores the file metadata, filename, and file content, returning a unique ID or an error.
	Save(metadata string, fileName string, fileData bytes.Buffer) (string, error)
}
