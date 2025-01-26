package app

import (
	"errors"
	"fmt"
	"github.com/kamchatkin/practicum-shortener/config"
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

// aliasProps Именнованные параметры для создания алиаса
type aliasProps struct {
	SourceURL string
	HTTPS     bool
	Host      string
}

// makeAlias
func makeAlias(props *aliasProps) (string, error) {
	var aliasKey string
	for i := range maxIterate {
		aliasKey = shortness()

		if _, ok := db[aliasKey]; !ok {
			break
		}

		i++
		if i == maxIterate {
			return "", errors.New("исчерпано максимальное количество попыток создания алиаса")
		}
	}

	db[aliasKey] = props.SourceURL

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
