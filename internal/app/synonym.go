package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"math/rand"
)

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
func makeAlias(ctx context.Context, props *aliasProps) (string, error) {
	var aliasKey string
	for i := range maxIterate {
		aliasKey = shortness()

		alias, err := storage.DB.Get(ctx, aliasKey)
		//  проблема взаимодействия с БД
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
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

	err := storage.DB.Set(ctx, aliasKey, props.SourceURL)
	if err != nil {
		return "", errors.Join(errors.New("не удалось записать в бд"), err)
	}

	proto := "http"
	if props.HTTPS {
		proto = "https"
	}

	host := props.Host

	cfg, _ := config.Config() // игнорируем ошибку потому что это ни на что не влияет
	if cfg.ShortHost != "" {
		proto = cfg.ShortHostURL.Scheme
		host = cfg.ShortHostURL.Host
	}

	return fmt.Sprintf("%s://%s/%s", proto, host, aliasKey), nil
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
