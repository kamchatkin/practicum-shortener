package app

import (
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSynonymHandler(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		method       string
		body         io.Reader
		expectedCode int
		expectedBody string
	}{
		{
			name:         "404 на POST",
			url:          "/synonym",
			method:       http.MethodPost,
			body:         strings.NewReader(""),
			expectedCode: http.StatusNotFound,
			expectedBody: "",
		},
		{
			name:         "400 неверные параметры",
			url:          "/",
			method:       http.MethodPost,
			body:         strings.NewReader("qwe"),
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name:         "Успешный тест",
			url:          "/",
			method:       http.MethodPost,
			body:         strings.NewReader("http://practicum.yandex.ru/"),
			expectedCode: http.StatusCreated,
			expectedBody: ".",
		},
	}

	storage.InitStorage()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, test.url, test.body)
			w := httptest.NewRecorder()
			SynonymHandler(w, r)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.expectedCode, res.StatusCode)

			if test.expectedBody == "ok" {
				_, err := url.ParseRequestURI(w.Body.String())
				assert.Nil(t, err, "В ответ должна быть получена правильная ссылка")
			}
		})
	}
}
