package file

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/models"
	"os"
	"sync"
	"time"
)

var memoryDB sync.Map

type FileStorage struct {
	isOpened bool
}

type dbRecord struct {
	Alias  string `json:"alias"`
	Source string `json:"source"`
}

// Set
func (f *FileStorage) Set(_ context.Context, key, value string) error {
	memoryDB.Store(key, value)

	return nil
}

// Get
func (f *FileStorage) Get(_ context.Context, key string) (models.Alias, error) {
	value, ok := memoryDB.Load(key)
	if !ok {
		return models.Alias{}, errors.New("not found")
	}

	return f.asAlias(key, value.(string)), nil
}
func (f *FileStorage) Incr() {}

// Open
func (f *FileStorage) Open() (bool, error) {
	f.isOpened = true
	cfg, _ := config.Config()

	file, err := os.OpenFile(cfg.DBFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		f.isOpened = false
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rec := dbRecord{}
		_ = json.Unmarshal([]byte(scanner.Text()), &rec)
		memoryDB.Store(rec.Alias, rec.Source)
	}

	return f.isOpened, nil
}

func (f *FileStorage) Opened() bool {
	return f.isOpened
}

// Close
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

	f.isOpened = false

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
