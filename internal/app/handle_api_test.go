package app

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestTo := httptest.NewRequest(http.MethodPost, "/api/shorten", tt.body)
			responseTo := httptest.NewRecorder()

			HandleAPI(responseTo, requestTo)

			response := responseTo.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expectedCode, response.StatusCode)
			if tt.expectedBody {
				responseBody := responseTo.Body.String()
				if assert.NotEmpty(t, responseBody) {
					resp := APIResponse{}
					err := json.NewDecoder(strings.NewReader(responseBody)).Decode(&resp)
					assert.NoError(t, err)

					err = validate.Struct(resp)
					assert.NoError(t, err)
				}
			}
		})
	}
}
