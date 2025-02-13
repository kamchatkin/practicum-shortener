package app

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"net/http"
	"time"
)

// RedirectHandler Поиск сокращения и переадресация
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	alias, err := storage.DB.Get(ctx, ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if alias.NotFound() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Location", alias.Source)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
