package app

import (
	"encoding/json"
	"fmt"
	"github.com/kamchatkin/practicum-shortener/internal/auth"
	"github.com/kamchatkin/practicum-shortener/internal/data"
	"github.com/kamchatkin/practicum-shortener/internal/logs"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"net/http"
)

func HandleDeleteAPIUserURLs(w http.ResponseWriter, r *http.Request) {

	var shortsToDelete []string
	json.NewDecoder(r.Body).Decode(&shortsToDelete)
	defer r.Body.Close()

	if len(shortsToDelete) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, _ := storage.NewStorage()
	userID := auth.GetUserIDFromCookie(r)
	logger := logs.NewLogger()

	go func() {
		err := data.UserBatchUpdate(db, userID, shortsToDelete)
		if err != nil {
			logger.Error(fmt.Errorf("HandleDeleteAPIUserURLs: failed to delete user: %w", err).Error())
		}
	}()

	w.WriteHeader(http.StatusAccepted)
}
