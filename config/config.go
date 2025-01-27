package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"net/url"
	"strings"
)

const DefaultAddr = ":8080"

const DefaultShortHost = ""

const DefaultDBStoragePath = "/tmp"

const DefaultDBDumpFileName = "db.txt"

type ConfigType struct {
	// Addr адрес для запуска этого приложения
	Addr string `env:"SERVER_ADDRESS" validate:"hostname_port"`

	// ShortHost подменный хост в сокращенном УРЛ
	ShortHost string `env:"BASE_URL" validate:"omitempty,http_url"`

	// DBStoragePath Путь хранения дампа БД
	DBStoragePath string `env:"FILE_STORAGE_PATH" validate:"dir"`

	ShortHostURL *url.URL

	// Возможность принудительного изменения полей. Исходно для тестов
	forceAddr      string
	forceShortHost string
}

func (cfg *ConfigType) DumpPath() string {
	return fmt.Sprintf("/%s/%s", strings.Trim(cfg.DBStoragePath, "/"), DefaultDBDumpFileName)
}

var hookAddr = ""

func HookAddr(val string) {
	hookAddr = val
}

var hookShortHost = ""

func HookShortHost(val string) {
	hookShortHost = val
}

// Config Конфигурация на момент запуска приложения
func Config() (*ConfigType, error) {
	cfg := &ConfigType{}
	parseArgs(cfg)
	if err := parseEnv(cfg); err != nil {
		return &ConfigType{}, err
	}

	ifTrue(&hookAddr, &cfg.Addr)
	ifTrue(&hookShortHost, &cfg.ShortHost)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(cfg)
	if err != nil {
		return &ConfigType{}, err
	}

	if cfg.ShortHost != "" {
		_url, err := url.Parse(cfg.ShortHost)
		if err != nil {
			return &ConfigType{}, err
		}
		cfg.ShortHostURL = _url
	}

	return cfg, nil
}

var parsedArgs = ConfigType{}

// parseArgs
func parseArgs(cfg *ConfigType) {
	if parsedArgs == (ConfigType{}) {
		flag.StringVar(&parsedArgs.Addr, "a", DefaultAddr, "Адрес запуска сервера. [HOST]:PORT")
		flag.StringVar(&parsedArgs.ShortHost, "b", DefaultShortHost, "Подменный УРЛ для сокращенного УРЛ. HOST[:PORT]")
		flag.StringVar(&parsedArgs.DBStoragePath, "f", DefaultDBStoragePath, "Путь до сохранения дампа БД")
		flag.Parse()
	}

	ifTrue(&parsedArgs.Addr, &cfg.Addr)
	ifTrue(&parsedArgs.ShortHost, &cfg.ShortHost)
	ifTrue(&parsedArgs.DBStoragePath, &cfg.DBStoragePath)
}

var parsedEnv = ConfigType{}

// parseEnv
func parseEnv(cfg *ConfigType) error {
	if parsedEnv == (ConfigType{}) {
		err := env.Parse(&parsedEnv)
		if err != nil {
			return err
		}
	}

	ifTrue(&parsedEnv.Addr, &cfg.Addr)
	ifTrue(&parsedEnv.ShortHost, &cfg.ShortHost)
	ifTrue(&parsedEnv.DBStoragePath, &cfg.DBStoragePath)

	return nil
}

func ifTrue(from, to *string) {
	if *from == "" {
		return
	}

	*to = *from
}
