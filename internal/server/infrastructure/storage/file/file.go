package file

import "bytes"

type Repository interface {
	Save(metadata string, fileName string, fileData bytes.Buffer) (string, error)
}
