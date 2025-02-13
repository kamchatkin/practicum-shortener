package memory

import (
	"context"
	"github.com/kamchatkin/practicum-shortener/internal/models"
	"sync"
	"time"
)

var memoryDB sync.Map

type MemStorage struct {
	isOpened bool
}

func (m *MemStorage) Set(_ context.Context, key, value string) error {
	memoryDB.Store(key, value)

	return nil
}

func (m *MemStorage) Get(_ context.Context, key string) (models.Alias, error) {
	value, ok := memoryDB.Load(key)
	if !ok {
		return models.Alias{}, nil
	}

	return m.asAlias(key, value.(string)), nil
}

func (m *MemStorage) Incr() {}
func (m *MemStorage) Open() (bool, error) {
	m.isOpened = true
	return m.isOpened, nil
}
func (m *MemStorage) Opened() bool {
	return m.isOpened
}
func (m *MemStorage) Close() error {
	m.isOpened = false

	return nil
}

func (m *MemStorage) asAlias(key, value string) models.Alias {
	return models.Alias{
		ID:        -1,
		Alias:     key,
		Source:    value,
		Quantity:  0,
		CreatedAt: time.Time{},
	}
}

func (m *MemStorage) Ping(_ context.Context) error {
	return nil
}
