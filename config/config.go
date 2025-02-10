package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"net/url"
)

const DefaultAddr = ":8080"

const DefaultShortHost = ""

const DefaultDBFilePath = "/tmp/1.db"

// ConfigType структура конфига
type ConfigType struct {
	// Addr адрес для запуска этого приложения
	Addr string `env:"SERVER_ADDRESS" validate:"hostname_port"`

	// ShortHost подменный хост в сокращенном УРЛ
	ShortHost    string `env:"BASE_URL" validate:"omitempty,http_url"`
	ShortHostURL *url.URL

	// DBFilePath Путь хранения дампа БД
	DBFilePath string `env:"FILE_STORAGE_PATH"`

	// Возможность принудительного изменения полей. Исходно для тестов
	forceAddr      string
	forceShortHost string
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
		flag.StringVar(&parsedArgs.DBFilePath, "f", DefaultDBFilePath, "Путь до сохранения дампа БД")
		flag.Parse()
	}

	ifTrue(&parsedArgs.Addr, &cfg.Addr)
	ifTrue(&parsedArgs.ShortHost, &cfg.ShortHost)
	ifTrue(&parsedArgs.DBFilePath, &cfg.DBFilePath)
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
	ifTrue(&parsedEnv.DBFilePath, &cfg.DBFilePath)

	return nil
}

func ifTrue(from, to *string) {
	if *from == "" {
		return
	}

	*to = *from
}
