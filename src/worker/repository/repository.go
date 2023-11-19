package repository

import (
	"fmt"
	"os"
)

type Repository struct {
	Path string
}

func NewRepository(path string) (*Repository, error) {
	tempFile, err := os.Create(path)

	if err != nil {
		return nil, err
	}

	return &Repository{
		Path: path,
	}, tempFile.Close()
}

func (m *Repository) Read() (int, error) {
	file, err := os.OpenFile(m.Path, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		return 0, err
	}

	var current int64

	if _, err = fmt.Fscanf(file, "%d", &current); err != nil && err.Error() != "EOF" {
		return 0, err
	}

	return int(current), file.Close()
}

func (m *Repository) Write(value int) error {
	file, err := os.OpenFile(m.Path, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(file, "%d", value); err != nil {
		return err
	}

	return file.Close()
}
