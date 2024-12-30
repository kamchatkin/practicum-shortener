package app

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestShortener(t *testing.T) {
	var baseURL = "http://localhost/?test"

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(baseURL))
	SynonymHandler(w, r)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	shortURL := w.Body.String()
	_, err := url.ParseRequestURI(shortURL)
	assert.Nil(t, err, "В ответ на сокращение ожидается URL")

	parts := strings.Split(shortURL, "/")

	w2 := httptest.NewRecorder()

	r2 := httptest.NewRequest(http.MethodGet, "/", nil)
	r2.SetPathValue("id", parts[len(parts)-1])
	RedirectHandler(w2, r2)

	resp2 := w2.Result()
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusTemporaryRedirect, resp2.StatusCode)
	assert.Equal(t, baseURL, w2.Header().Get("Location"))
}
