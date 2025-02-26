package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kamchatkin/practicum-shortener/internal/auth"
	"github.com/kamchatkin/practicum-shortener/internal/data"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

var cookieWithAliases = ""
var cookieWithoutAliases = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzE2MzEwNDgsIlVzZXJJRCI6NDl9.f1s6RRuuqx72neVxpxDb40vGpAdjg5sv5l2X8FVECCo"

func TestHandleAPIUserURLs(t *testing.T) {
	tests := []struct {
		name        string
		validCookie bool
		expCode     int
		expBody     bool
		expHeader   bool
		aliases     bool
	}{
		{
			name:        "Неизвестный пользователь",
			validCookie: false,
			expCode:     http.StatusUnauthorized,
		},
		{
			name:        "ОК. Один сокрщенный УРЛ",
			validCookie: true,
			expCode:     http.StatusOK,
			expBody:     true,
			expHeader:   true,
			aliases:     true,
		},
		{
			name:        "OK. No content",
			validCookie: true,
			expCode:     http.StatusNoContent,
			expBody:     false,
			expHeader:   false,
			aliases:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cookieStr := "---"
			userID := int64(-1)

			if tt.validCookie {
				cookieStr, _ = auth.BuildJWTString()
				userID = auth.GetUserID(cookieStr)
				fmt.Printf("userID: %d\n", userID)
			}

			if tt.aliases {
				db, _ := storage.NewStorage()
				err := data.Set(context.TODO(), db, shortness(), fmt.Sprintf("https://asdfghj:9000/difnisdf?sdfjni1=%d", rand.Int64()), userID)
				assert.NoError(t, err)
			}

			r := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			r.AddCookie(&http.Cookie{
				Name:  auth.CookineName,
				Value: cookieStr,
			})

			w := httptest.NewRecorder()

			HandleAPIUserURLs(w, r)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.expCode, result.StatusCode)

			if tt.expBody {
				var urlsResp []*UserURLsResponse
				json.NewDecoder(result.Body).Decode(&urlsResp)
				assert.GreaterOrEqual(t, len(urlsResp), 1) // один и более сокращенный URL
			}

			if tt.expHeader {
				assert.Contains(t, result.Header.Get("Content-Type"), "application/json")
			}
		})
	}
}
