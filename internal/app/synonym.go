package app

import "math/rand"

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
