package storage

import (
	"context"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/logs"
	"github.com/kamchatkin/practicum-shortener/internal/models"
	fileStorage "github.com/kamchatkin/practicum-shortener/internal/storage/file"
	memoryStorage "github.com/kamchatkin/practicum-shortener/internal/storage/memory"
	pgStorage "github.com/kamchatkin/practicum-shortener/internal/storage/pg"
)

// DB Публичный доступ к хранилищу
var DB Storage

// Storage интерфейс описывающий требования к хранилищу
type Storage interface {
	// Set запись в хранилище
	Set(ctx context.Context, key string, value string) error

	// Get Получает данные из хранилища по ключу
	Get(ctx context.Context, key string) (models.Alias, error)

	// Incr Инкриминирует счетчик переходов по сокращению
	Incr()

	// Open Открывает хранилище. При загрузке приложения
	Open() (bool, error)

	// Opened Открыто ли хранилище? @deprecated
	Opened() bool

	// Close Закрытие хранилища
	Close() error

	// Ping Тестовый запрос в хранилище для проверки работоспособности
	Ping(ctx context.Context) error
}

func InitStorage() {
	cfg, _ := config.Config()

	logger := logs.NewLogger()

	if DB != nil {
		logger.Info("Storage already initialized")
		return
	}

	if cfg.DatabaseDsn != "" {
		logger.Info("Выбрана БД: postgresql")
		DB = &(pgStorage.PostgresStorage{})
		return
	}

	if cfg.DBFilePath != "" && cfg.DBFilePath != config.DefaultDBFilePath {
		logger.Info("Выбрана БД: файловое хранилище")
		DB = &(fileStorage.FileStorage{})
		return
	}

	logger.Info("Выбрана БД: в памяти приложения (до перезагрузки)")
	DB = &(memoryStorage.MemStorage{})
}
