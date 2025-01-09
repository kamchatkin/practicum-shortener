package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = []string{"app", "-a", "...."}
	assert.Panics(t, main)
}
