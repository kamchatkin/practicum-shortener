package memory

import (
	"context"
	"errors"
	"fmt"
	"github.com/kamchatkin/practicum-shortener/internal/models"
	"math/rand"
	"sync"
	"time"
)

var memoryDB sync.Map
var mems *MemStorage

var (
	linksMU sync.RWMutex
	links   map[int64][]string
)

type UniqError error

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

func (m *MemStorage) Set(_ context.Context, key, value string, userID int64) error {
	uniqErr := false
	memoryDB.Range(func(mKey, mValue any) bool {
		if value == mValue {
			uniqErr = true
			return false
		}

		return true
	})

	if uniqErr {
		return UniqError(fmt.Errorf("duplicate value"))
	}

	memoryDB.Store(key, value)
	linksMU.Lock()
	if _, ok := links[userID]; !ok {
		links[userID] = []string{}
	}
	links[userID] = append(links[userID], key)
	linksMU.Unlock()

	return nil
}

func (m *MemStorage) SetBatch(_ context.Context, item map[string]string, userID int64) error {
	linksMU.Lock()
	for key, value := range item {
		memoryDB.Store(key, value)

		if _, ok := links[userID]; !ok {
			links[userID] = []string{}
		}
		links[userID] = append(links[userID], key)
	}
	linksMU.Unlock()

	return nil
}

func (m *MemStorage) Get(_ context.Context, key string) (models.Alias, error) {
	value, ok := memoryDB.Load(key)
	if !ok {
		return models.Alias{}, nil
	}

	return m.asAlias(key, value.(string)), nil
}

func (m *MemStorage) GetBySource(_ context.Context, source string) (models.Alias, error) {
	shortKey := ""
	memoryDB.Range(func(mKey, mValue any) bool {
		if mValue == source {
			shortKey = mKey.(string)
		}

		return true
	})

	return m.asAlias(shortKey, source), nil
}

func (m *MemStorage) RegisterUser(_ context.Context) (int64, error) {
	return rand.Int63(), nil
}

func (m *MemStorage) UserAliases(_ context.Context, userID int64) ([]*models.Alias, error) {
	var aliases []*models.Alias

	if userID < 0 {
		return aliases, errors.New("invalid userID")
	}

	linksMU.RLock()
	defer linksMU.RUnlock()

	userShorts, ok := links[userID]
	if !ok {
		return aliases, nil
	}

	for _, uShort := range userShorts {
		value, _ := memoryDB.Load(uShort)
		alias := m.asAlias(uShort, value.(string))
		aliases = append(aliases, &alias)
	}

	return aliases, nil
}

func (m *MemStorage) Incr() {}
func (m *MemStorage) Open() error {
	links = map[int64][]string{}
	links[-1] = []string{"qwerty"}
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

func (m *MemStorage) IsUniqError(err error) bool {
	return errors.Is(err, UniqError(err))
}

func (m *MemStorage) UserBatchUpdate(_ context.Context, _ chan string, _ int64) error {
	return nil
}
