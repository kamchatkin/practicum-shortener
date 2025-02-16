package memory

import (
	"context"
	"fmt"
	"github.com/kamchatkin/practicum-shortener/internal/models"
	"sync"
	"time"
)

var memoryDB sync.Map
var mems *MemStorage

func init() {
	memoryDB.Store("qwerty", "https://ya.ru/")
}

type MemStorage struct{}

func NewMemStorage() (*MemStorage, error) {
	if mems != nil {
		return mems, nil
	}

	mems = &MemStorage{}
	err := mems.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open memory storage: %w", err)
	}

	return mems, nil
}

func (m *MemStorage) Set(_ context.Context, key, value string) error {
	memoryDB.Store(key, value)

	return nil
}

func (m *MemStorage) SetBatch(_ context.Context, item map[string]string) error {
	for key, value := range item {
		memoryDB.Store(key, value)
	}

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
func (m *MemStorage) Open() error {
	return nil
}
func (m *MemStorage) Close() error {
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
