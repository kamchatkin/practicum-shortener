package app

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleAPIBatch(t *testing.T) {

	ids := "123|456"

	requestJSON := `[ {"correlation_id": "123", "original_url": "http://ya.ru/"},
					  {"correlation_id": "456", "original_url": "http://ya.ru/?1"}]`

	tests := []struct {
		name          string
		expCode       int
		expBody       string
		expJSONLength int
	}{
		{
			name:          "success",
			expCode:       http.StatusCreated,
			expJSONLength: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(requestJSON))
			w := httptest.NewRecorder()

			HandleAPIBatch(w, r)

			response := w.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.expCode, response.StatusCode)

			var respJSON []BatchItemTo
			err := json.NewDecoder(response.Body).Decode(&respJSON)
			assert.NoError(t, err)

			assert.Equal(t, tt.expJSONLength, len(respJSON))

			for _, alias := range respJSON {
				assert.NotEmpty(t, alias.CorrelationID)
				assert.NotEmpty(t, alias.ShortURL)

				assert.Contains(t, ids, alias.CorrelationID)
			}
		})
	}
}
