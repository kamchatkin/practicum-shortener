package app

import (
	"context"
	"encoding/json"
	"github.com/kamchatkin/practicum-shortener/internal/auth"
	"github.com/kamchatkin/practicum-shortener/internal/data"
	"github.com/kamchatkin/practicum-shortener/internal/logs"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"go.uber.org/zap"
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

	logger := logs.NewLogger()

	userID := auth.GetUserIDFromCookie(r)
	logger.Info("userID", zap.Int64("userID", userID))
	//if userID < 1 {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	db, err := storage.NewStorage()
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	aliases, err := data.UserAliases(ctx, db, userID)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(aliases) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	aliasProperties := aliasProps{
		SourceURL: "",
		HTTPS:     r.TLS != nil,
		Host:      r.Host,
	}
	var responseAliases []*UserURLsResponse
	for _, alias := range aliases {
		responseAliases = append(responseAliases, &UserURLsResponse{
			ShortURL:    getShortURL(alias.Alias, &aliasProperties),
			OriginalURL: alias.Source,
		})
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(responseAliases)
}
