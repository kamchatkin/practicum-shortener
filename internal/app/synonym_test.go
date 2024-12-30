package app

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

// Test_shortness
func Test_shortness(t *testing.T) {
	alias := shortness()
	assert.NotEmpty(t, alias, "Ожидается не пустая строка")
	assert.Regexp(t, regexp.MustCompile(`^[a-zA-Z0-9]+$`), alias, "Сформированная строка должна состоять из латиницы (строчной и заглавной) и цифр.")
	assert.Len(t, []rune(alias), LENGTH)
}

// Test_randInt
func Test_randInt(t *testing.T) {
	minValue := 1
	maxValue := 5
	val := randInt(minValue, maxValue)
	assert.True(t, minValue <= val && val <= maxValue,
		fmt.Sprintf("Сгенерированное значение должно быть не меньше минимального (%d) и не больше максимального (%d). Получено %d", minValue, maxValue, val))
}
