package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var baseArgs = []string{"shortener.exe", "-a", validA, "-b", validB}

var validA = ":8081" // @todo потенциальная проблема. :8888 проходит проверку в тесте, но не в конфиге. Это как так?
var validA2 = ":8082"
var invalidA = "Невозможно :D"

var validB = "https://vk.com/"
var validB2 = "https://vk2.com/"
var invalidB = "vk.com"

// TestParseArgs Тест аргументов
func TestParseArgs(t *testing.T) {
	os.Args = baseArgs
	cfg, err := Config()
	assert.NoError(t, err)
	assert.Equal(t, validA, cfg.Addr)
	assert.Equal(t, validB, cfg.ShortHost)
}

// TestParseEnv Тест переменных окружения
func TestParseEnv(t *testing.T) {
	parsedEnv = ConfigType{}
	_ = os.Setenv("SERVER_ADDRESS", validA2)
	_ = os.Setenv("BASE_URL", validB2)
	cfg, err := Config()
	assert.NoError(t, err)
	assert.Equal(t, validA2, cfg.Addr)
	assert.Equal(t, validB2, cfg.ShortHost)
}

// TestParseError Ошибка валидации
func TestParseError(t *testing.T) {
	parsedEnv = ConfigType{}
	_ = os.Setenv("SERVER_ADDRESS", invalidA)
	_ = os.Setenv("BASE_URL", invalidB)
	cfg, err := Config()
	assert.Error(t, err)
	assert.Equal(t, "", cfg.Addr)
	assert.Equal(t, "", cfg.ShortHost)
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
	cc := ConfigType{Addr: ":8080"}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(cc)

	assert.NoError(t, err)
}
