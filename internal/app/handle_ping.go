package app

import (
	"context"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"net/http"
	"time"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	db, err := storage.NewStorage()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	err = storage.Ping(ctx, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
