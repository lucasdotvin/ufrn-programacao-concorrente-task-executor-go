package repository

import (
	"encoding/binary"
	"os"
)

type Repository struct {
	Path             string
	LastWrittenValue int
}

func NewRepository(path string) (*Repository, error) {
	tempFile, err := os.Create(path)

	if err != nil {
		return nil, err
	}

	return &Repository{
		Path:             path,
		LastWrittenValue: 0,
	}, tempFile.Close()
}

func (m *Repository) Read() (int, error) {
	return m.LastWrittenValue, nil
}

func (m *Repository) Write(value int) error {
	file, err := os.OpenFile(m.Path, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	content := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(content, int64(value))

	if _, err = file.Write(content); err != nil {
		return err
	}

	m.LastWrittenValue = value

	return file.Close()
}
