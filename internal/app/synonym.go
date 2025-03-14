package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/data"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"math/rand"
)

var ErrUniq = errors.New("errUniq")

var words []rune
var wordsQuantity = 0

func init() {
	for a := 'a'; a <= 'z'; a++ {
		words = append(words, a)
	}

	for a := 'A'; a <= 'Z'; a++ {
		words = append(words, a)
	}

	for a := '0'; a <= '9'; a++ {
		words = append(words, a)
	}
	wordsQuantity = len(words) - 1
}

// LENGTH длина алиаса для сокращения
const LENGTH = 5

// aliasProps Именованные параметры для создания алиаса
type aliasProps struct {
	SourceURL string
	HTTPS     bool
	Host      string
}

// makeAlias
func makeAlias(ctx context.Context, db *storage.Storage, props *aliasProps, userID int64) (string, error) {
	aliasKey, err := getShortCode(ctx, db)
	if err != nil {
		return "", fmt.Errorf("could not get short code for alias: %w", err)
	}

	err = data.Set(ctx, db, aliasKey, props.SourceURL, userID)
	if err != nil {

		if (*db).IsUniqError(err) {
			origShortURL, err := SearchOriginalALias(ctx, db, props.SourceURL, props)
			if err != nil {
				return "", err
			}

			return origShortURL, ErrUniq
		}

		return "", errors.Join(errors.New("не удалось записать в бд"), err)
	}

	return getShortURL(aliasKey, props), nil
}

// SearchOriginalALias
func SearchOriginalALias(ctx context.Context, db *storage.Storage, sourceURL string, props *aliasProps) (string, error) {
	alias, err := data.GetBySource(ctx, db, sourceURL)
	if err != nil {
		return "", err
	}

	return getShortURL(alias.Alias, props), nil
}

// @todo можно облегчить за счет контролируемого создания уникального кода
func getShortCode(ctx context.Context, db *storage.Storage) (string, error) {
	var aliasKey string
	for i := range maxIterate {
		aliasKey = shortness()

		alias, err := data.Get(ctx, db, aliasKey)
		// проблема взаимодействия с БД
		if err != nil {
			return "", fmt.Errorf("failed to get alias for %s: %w", aliasKey, err)
		}

		if alias.NotFound() {
			break
		}

		i++
		if i == maxIterate {
			return "", errors.New("исчерпано максимальное количество попыток создания алиаса")
		}
	}

	return aliasKey, nil
}

func getShortURL(shortCode string, prop *aliasProps) string {
	proto := "http"
	if prop.HTTPS {
		proto = "https"
	}

	host := prop.Host

	cfg, _ := config.Config() // игнорируем ошибку потому что это ни на что не влияет
	if cfg.ShortHost != "" {
		proto = cfg.ShortHostURL.Scheme
		host = cfg.ShortHostURL.Host
	}

	return fmt.Sprintf("%s://%s/%s", proto, host, shortCode)
}

// shortness
func shortness() string {
	var str []rune
	for range LENGTH {
		str = append(str, words[randInt(0, wordsQuantity)])
	}

	return string(str)
}

// randInt
func randInt(a, b int) int {
	return a + rand.Intn(b-a+1)
}
