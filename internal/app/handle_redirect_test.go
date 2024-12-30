package app

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Тестирование переадресации
func TestRedirect(t *testing.T) {
	redirectTestCases := []struct {
		url          string
		method       string
		expectedCode int
		expectedBody string
		codeMsg      string
		bodyMsg      string
	}{
		{
			url:          "1",
			method:       http.MethodGet,
			expectedCode: http.StatusNotFound,
			expectedBody: "",
			codeMsg:      fmt.Sprintf("GET запрос несуществующего ID сокращения должен возвращать %d", http.StatusNotFound),
			bodyMsg:      "В ответ ожидается пустое тело ответа",
		},
		{
			url:          "qwerty",
			method:       http.MethodGet,
			expectedCode: http.StatusTemporaryRedirect,
			expectedBody: "",
			codeMsg:      fmt.Sprintf("GET запрос должен возвращать код временной переадресации %d", http.StatusTemporaryRedirect),
			bodyMsg:      "",
		},
	}

	for _, tc := range redirectTestCases {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/", nil)
			r.SetPathValue("id", tc.url)
			w := httptest.NewRecorder()

			RedirectHandler(w, r)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode, tc.codeMsg)
			if tc.expectedBody == "" {
				assert.Empty(t, w.Body.String(), tc.bodyMsg)
			} else {
				assert.Equal(t, tc.expectedBody, w.Body.String(), tc.bodyMsg)
			}
		})
	}

}
