package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func baseURL() string {
	return "http://localhost/?q=" + shortness()
}

func TestShortener(t *testing.T) {
	_base := baseURL()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(_base))
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
	//r2.SetPathValue("id", parts[len(parts)-1])
	r2ctx := chi.NewRouteContext()
	r2ctx.URLParams.Add("id", parts[len(parts)-1])
	r2 = r2.WithContext(context.WithValue(r2.Context(), chi.RouteCtxKey, r2ctx))

	RedirectHandler(w2, r2)

	resp2 := w2.Result()
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusTemporaryRedirect, resp2.StatusCode)
	assert.Equal(t, _base, w2.Header().Get("Location"))
}

func TestShortenerWithConfig(t *testing.T) {
	shortHost := "http://ya.ru"
	_base := baseURL()

	config.HookShortHost(shortHost)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(_base))
	SynonymHandler(w, r)

	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	shortURL := w.Body.String()
	_url, err := url.Parse(shortURL)

	assert.Nil(t, err, fmt.Sprintf("В ответ на сокращение ожидается URL, получено: %s", err))
	assert.Equal(t, shortHost, fmt.Sprintf("%s://%s", _url.Scheme, _url.Host), "Значение shortHost конфигурации должно учитываться в короткой ссылке")
}
