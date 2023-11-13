package repository

import (
	"encoding/binary"
	"os"
	"sync"
)

type Repository struct {
	Path    string
	FileMux *sync.RWMutex
}

func NewRepository(path string, fileMux *sync.RWMutex) *Repository {
	return &Repository{
		Path:    path,
		FileMux: fileMux,
	}
}

func (m *Repository) Read() (int, error) {
	m.FileMux.RLock()
	defer m.FileMux.RUnlock()

	file, err := os.OpenFile(m.Path, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		return 0, err
	}

	var current int64
	err = binary.Read(file, binary.LittleEndian, &current)

	if err != nil && err.Error() != "EOF" {
		return 0, err
	}

	return int(current), file.Close()
}

func (m *Repository) Write(value int) error {
	m.FileMux.Lock()
	defer m.FileMux.Unlock()

	file, err := os.OpenFile(m.Path, os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	if err := binary.Write(file, binary.LittleEndian, int64(value)); err != nil {
		return err
	}

	return file.Close()
}
