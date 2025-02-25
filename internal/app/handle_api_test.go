package app

import (
	"encoding/json"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// @todo нужно использовать моку для теста. Сейчас не работает

func TestHandleAPI(t *testing.T) {
	tests := []struct {
		name         string
		body         io.Reader
		expectedBody bool
		expectedCode int
	}{
		{
			name:         "Успешное создание сокращения",
			body:         strings.NewReader("{\"url\":\"https://practicum.yandex.ru/\"}"),
			expectedBody: true,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Ошибка. Не валидный JSON в запросе",
			body:         strings.NewReader("{\"url\"\"https://www.yandex.ru/\"}"),
			expectedBody: false,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Ошибка. Не валидный URL в запросе",
			body:         strings.NewReader("{\"url\":\"http2//:www.yandex.ru/\"}"),
			expectedBody: false,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "iter13. Дубликат длинного URL",
			body:         strings.NewReader("{\"url\":\"" + config.DefaultSource + "\"}"),
			expectedBody: true,
			expectedCode: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestTo := httptest.NewRequest(http.MethodPost, "/api/shorten", tt.body)
			requestTo.Host = "localhost"
			responseTo := httptest.NewRecorder()

			HandleAPI(responseTo, requestTo)

			response := responseTo.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedCode, response.StatusCode)
			if tt.expectedBody {
				resp := APIResponse{}
				err := json.NewDecoder(response.Body).Decode(&resp)
				assert.NoError(t, err)

				err = validate.Struct(resp)
				assert.NoError(t, err)
			}
		})
	}
}
