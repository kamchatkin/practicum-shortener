package file

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/logs"
	"github.com/kamchatkin/practicum-shortener/internal/models"
	"math/rand"
	"os"
	"sync"
	"time"
)

var memoryDB = &sync.Map{}
var fStorage *FileStorage

var (
	linksMU sync.RWMutex
	links   map[int64][]string
)

type UniqError error

func init() {
	memoryDB.Store("qwerty", "https://ya.ru/")
}

type FileStorage struct{}

func NewFileStorage() (*FileStorage, error) {
	if fStorage != nil {
		return fStorage, nil
	}

	fs := &FileStorage{}
	logger := logs.NewLogger()
	err := fs.Open()
	if err != nil {
		err = fmt.Errorf("could not open file storage: %w", err)
		logger.Error(err.Error())

		return nil, err
	}

	return fs, nil
}

type dbRecord struct {
	Alias  string `json:"alias"`
	Source string `json:"source"`
}

// Set
func (f *FileStorage) Set(_ context.Context, key, value string, userID int64) error {
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

	memoryDB.Store(key, value)
	linksMU.Lock()
	if _, ok := links[userID]; !ok {
		links[userID] = []string{}
	}
	links[userID] = append(links[userID], key)
	linksMU.Unlock()

	return nil
}

func (f *FileStorage) SetBatch(_ context.Context, item map[string]string, userID int64) error {
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

// Get
func (f *FileStorage) Get(_ context.Context, key string) (models.Alias, error) {
	value, ok := memoryDB.Load(key)
	if !ok {
		return models.Alias{}, nil
	}

	return f.asAlias(key, value.(string)), nil
}

func (f *FileStorage) GetBySource(_ context.Context, source string) (models.Alias, error) {
	shortKey := ""
	memoryDB.Range(func(mKey, mValue any) bool {
		if mValue == source {
			shortKey = mKey.(string)
		}

		return true
	})

	return f.asAlias(shortKey, source), nil
}

func (f *FileStorage) RegisterUser(_ context.Context) (int64, error) {
	return rand.Int63(), nil
}

func (f *FileStorage) UserAliases(_ context.Context, userID int64) ([]*models.Alias, error) {
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
		alias := f.asAlias(uShort, value.(string))
		aliases = append(aliases, &alias)
	}

	return aliases, nil
}

func (f *FileStorage) Incr() {}

// Open чтение с диска
func (f *FileStorage) Open() error {
	cfg, _ := config.Config()

	file, err := os.OpenFile(cfg.DBFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	linksMU.Lock()
	defer linksMU.Unlock()

	links = map[int64][]string{}
	links[-1] = []string{"qwerty"}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rec := dbRecord{}
		_ = json.Unmarshal([]byte(scanner.Text()), &rec)
		memoryDB.Store(rec.Alias, rec.Source)
		links[-1] = append(links[-1], rec.Alias)
	}

	return nil
}

// Close Сохранение на диск
func (f *FileStorage) Close() error {
	cfg, err := config.Config()
	if err != nil {
		return err
	}

	err = os.WriteFile(cfg.DBFilePath, []byte{}, 0666)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(cfg.DBFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	memoryDB.Range(func(k, v interface{}) bool {
		rec := dbRecord{
			Alias:  k.(string),
			Source: v.(string),
		}
		_ = json.NewEncoder(writer).Encode(&rec)

		return true
	})

	return nil
}

func (f *FileStorage) asAlias(key, value string) models.Alias {
	return models.Alias{
		ID:        -1,
		Alias:     key,
		Source:    value,
		Quantity:  0,
		CreatedAt: time.Time{},
	}
}

func (f *FileStorage) Ping(_ context.Context) error {
	return nil
}

func (f *FileStorage) IsUniqError(err error) bool {
	return errors.Is(err, UniqError(err))
}
