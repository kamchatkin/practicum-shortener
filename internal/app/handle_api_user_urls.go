package app

import (
	"context"
	"encoding/json"
	"github.com/kamchatkin/practicum-shortener/internal/auth"
	"github.com/kamchatkin/practicum-shortener/internal/data"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"net/http"
	"time"
)

type UserURLsResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func HandleAPIUserURLs(w http.ResponseWriter, r *http.Request) {
	// Иметь хендлер GET /api/user/urls,
	// который сможет вернуть пользователю
	// все когда-либо сокращённые им URL в формате:

	userCookie, err := r.Cookie(auth.CookineName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID := auth.GetUserID(userCookie.Value)
	if userID < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	db, err := storage.NewStorage()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	aliases, err := data.UserAliases(ctx, db, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(aliases) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var responseAliases []*UserURLsResponse
	for _, alias := range aliases {
		responseAliases = append(responseAliases, &UserURLsResponse{
			ShortURL:    alias.Alias,
			OriginalURL: alias.Source,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(responseAliases)
}
