package storage

import (
	"context"
	"fmt"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/logs"
	"github.com/kamchatkin/practicum-shortener/internal/models"
	fileStorage "github.com/kamchatkin/practicum-shortener/internal/storage/file"
	memoryStorage "github.com/kamchatkin/practicum-shortener/internal/storage/memory"
	pgStorage "github.com/kamchatkin/practicum-shortener/internal/storage/pg"
)

// DB Публичный доступ к хранилищу
var dbRef Storage

// Storage интерфейс описывающий требования к хранилищу
type Storage interface {
	// Set запись в хранилище
	Set(ctx context.Context, key string, value string, userID int64) error

	SetBatch(ctx context.Context, item map[string]string, userID int64) error

	// Get Получает данные из хранилища по ключу
	Get(ctx context.Context, key string) (models.Alias, error)

	// GetBySource поиск по длинной ссылке
	GetBySource(ctx context.Context, key string) (models.Alias, error)

	// Incr Инкриминирует счетчик переходов по сокращению
	Incr()

	// Open Открывает хранилище. При загрузке приложения
	Open() error

	// Close Закрытие хранилища
	Close() error

	// Ping Тестовый запрос в хранилище для проверки работоспособности
	Ping(ctx context.Context) error

	IsUniqError(err error) bool

	RegisterUser(ctx context.Context) (int64, error)

	UserAliases(ctx context.Context, userID int64) ([]*models.Alias, error)

	UserBatchUpdate(ctx context.Context, shortsCh chan string, userID int64) error
}

func NewStorage() (*Storage, error) {
	cfg, _ := config.Config()

	if dbRef != nil {
		return &dbRef, nil
	}

	logger := logs.NewLogger()
	var err error

	if cfg.DatabaseDsn != "" && !cfg.TestENV {
		logger.Info("Выбрана БД: postgresql")
		dbRef, err = pgStorage.NewPostgresStorage()
		if err != nil {
			return nil, fmt.Errorf("could not connect to postgresql: %w", err)
		}
	}

	if cfg.DBFilePath != "" && !cfg.TestENV {
		logger.Info("Выбрана БД: файловое хранилище")
		dbRef, err = fileStorage.NewFileStorage()
		if err != nil {
			return nil, fmt.Errorf("could not connect to fileStorage: %w", err)
		}
	}

	if dbRef == nil || cfg.TestENV {
		logger.Info("Выбрана БД: в памяти приложения (до перезагрузки)")
		dbRef, err = memoryStorage.NewMemStorage()
		if err != nil {
			return nil, fmt.Errorf("could not connect to memoryStorage: %w", err)
		}
	}

	return &dbRef, nil
}

func Close() {
	if dbRef != nil {
		_ = dbRef.Close()
	}
}

func Ping(ctx context.Context, db *Storage) error {
	return (*db).Ping(ctx)
}
