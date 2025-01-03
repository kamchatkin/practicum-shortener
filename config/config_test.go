package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var validA = ":8888"
var validA2 = "Невозможно :D"

var validB = "https://vk.com/"
var invalidB = "vk.com"

func TestParse(t *testing.T) {
	os.Args = []string{"shortener.exe", "-a", validA, "-b", validB}
	Parse()

	assert.Equal(t, validA, Config.Addr)
	assert.Equal(t, validB, Config.ShortHost)

	err := os.Setenv("SERVER_ADDRESS", validA2)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("BASE_URL", validB)
	if err != nil {
		t.Fatal(err)
	}

	parseEnv()
	assert.Equal(t, validA2, Config.Addr)
	assert.Equal(t, validB, Config.ShortHost)

	assert.Panics(t, func() {
		confShortHost(invalidB)
	})
}
