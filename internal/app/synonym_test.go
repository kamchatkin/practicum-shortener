package app

import (
	"context"
	"fmt"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

// Test_shortness
func Test_SynonymShortness(t *testing.T) {
	alias := shortness()
	assert.NotEmpty(t, alias, "Ожидается не пустая строка")
	assert.Regexp(t, regexp.MustCompile(`^[a-zA-Z0-9]+$`), alias, "Сформированная строка должна состоять из латиницы (строчной и заглавной) и цифр.")
	assert.Len(t, []rune(alias), LENGTH)
}

// Test_randInt
func Test_SynonymRandInt(t *testing.T) {
	minValue := 1
	maxValue := 5
	val := randInt(minValue, maxValue)
	assert.True(t, minValue <= val && val <= maxValue,
		fmt.Sprintf("Сгенерированное значение должно быть не меньше минимального (%d) и не больше максимального (%d). Получено %d", minValue, maxValue, val))
}

func Test_SynonymGetShortCode(t *testing.T) {
	db, _ := storage.NewStorage()
	aliasStr, err := getShortCode(context.TODO(), db)
	assert.NoError(t, err)
	assert.NotEmpty(t, aliasStr)
}

func Test_SynonymGetShortURL(t *testing.T) {
	config.HookShortHost("")
	cfg, _ := config.Config()

	url := getShortURL(cfg.Addr, &aliasProps{
		HTTPS: false,
		Host:  "jj",
	})

	assert.NotEmpty(t, url)
	assert.Equal(t, "http://jj/:8080", url)
}

func Test_SynonymMakeAlias(t *testing.T) {
	config.HookShortHost("")
	db, _ := storage.NewStorage()
	alias, err := makeAlias(context.Background(), db, &aliasProps{
		SourceURL: "http://jj.ru/",
		HTTPS:     false,
		Host:      "jj:30001",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, alias)
	assert.Contains(t, alias, "http://jj:30001")
}
