package app

import (
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePing(t *testing.T) {
	r := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	HandlePing(w, r)

	response := w.Result()
	defer response.Body.Close()

	_, err := storage.NewStorage()
	if err != nil {
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	} else {
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}

}
