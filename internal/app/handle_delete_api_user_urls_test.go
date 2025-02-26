package app

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleDeleteAPIUserURLs(t *testing.T) {

	tests := []struct {
		name    string
		expCode int
		body    io.Reader
	}{
		{
			name:    "success",
			expCode: http.StatusAccepted,
			body:    strings.NewReader(`["qwerty", "undefined", "jsbfgij4"]`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodDelete, "/", tt.body)
			w := httptest.NewRecorder()

			HandleDeleteAPIUserURLs(w, r)

			response := w.Result()

			assert.Equal(t, tt.expCode, response.StatusCode)
		})
	}
}
