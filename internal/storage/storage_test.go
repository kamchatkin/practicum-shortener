package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStorage(t *testing.T) {
	st, err := NewStorage()
	assert.NoError(t, err)
	assert.NotNil(t, st)
}
