package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var baseArgs = []string{"shortener.exe", "-a", validA, "-b", validB, "-f", validF, "-d", validD}

var validA = ":8081"
var validA2 = ":8082"
var invalidA = "Невозможно :D"

var validB = "https://vk.com/"
var validB2 = "https://vk2.com/"
var invalidB = "vk.com"

var validF = "/tmp"
var validF2 = "/dev/null"
var invalidF = "/tmp/.a"

var validD = "postgres://username:password@localhost:5432/database_name1"
var validD2 = "postgres://username:password@localhost:5432/database_name2"
var invalidD = "postgresusernamepasswordlocalhost5432database_name" // postgres://username:password@localhost:5432/database_name

// TestParseArgs Тест аргументов
func TestParseArgs(t *testing.T) {
	os.Args = baseArgs
	cfg, err := Config()
	assert.NoError(t, err)
	assert.Equal(t, validA, cfg.Addr)
	assert.Equal(t, validB, cfg.ShortHost)
	assert.Equal(t, validF, cfg.DBFilePath)
	assert.Equal(t, validD, cfg.DatabaseDsn)
}

// TestParseEnv Тест переменных окружения
func TestParseEnv(t *testing.T) {
	parsedEnv = ConfigType{}
	_ = os.Setenv("SERVER_ADDRESS", validA2)
	_ = os.Setenv("BASE_URL", validB2)
	_ = os.Setenv("FILE_STORAGE_PATH", validF2)
	_ = os.Setenv("DATABASE_DSN", validD2)
	cfg, err := Config()
	assert.NoError(t, err)
	assert.Equal(t, validA2, cfg.Addr)
	assert.Equal(t, validB2, cfg.ShortHost)
	assert.Equal(t, validF2, cfg.DBFilePath)
	assert.Equal(t, validD2, cfg.DatabaseDsn)
}

// TestParseError Ошибка валидации
func TestParseError(t *testing.T) {
	parsedEnv = ConfigType{}
	_ = os.Setenv("SERVER_ADDRESS", invalidA)
	_ = os.Setenv("BASE_URL", invalidB)
	_ = os.Setenv("FILE_STORAGE_PATH", invalidF)
	_ = os.Setenv("DATABASE_DSN", invalidD)
	cfg, err := Config()
	assert.Error(t, err)
	assert.Equal(t, "", cfg.Addr)
	assert.Equal(t, "", cfg.ShortHost)
	assert.Equal(t, "", cfg.DBFilePath)
	assert.Equal(t, "", cfg.DatabaseDsn)
}

func TestIfTrue(t *testing.T) {
	type Aa struct {
		Aa string
	}

	a := Aa{Aa: "Aa"}
	b := Aa{Aa: "Bb"}

	ifTrue(&a.Aa, &b.Aa)
	assert.Equal(t, a.Aa, b.Aa)
}

func TestValidate(t *testing.T) {
	cc := ConfigType{Addr: ":8080", DBFilePath: "/tmp"}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(cc)

	assert.NoError(t, err)
}
