package app

import (
	"context"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"github.com/kamchatkin/practicum-shortener/internal/storage/pg"
	"net/http"
	"time"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	switch storage.DB.(type) {
	case *pg.PostgresStorage:

		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		if err := storage.DB.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
