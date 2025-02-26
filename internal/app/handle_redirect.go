package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kamchatkin/practicum-shortener/internal/data"
	"github.com/kamchatkin/practicum-shortener/internal/logs"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"net/http"
	"time"
)

// RedirectHandler Поиск сокращения и переадресация
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	logger := logs.NewLogger()

	db, err := storage.NewStorage()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	alias, err := data.Get(ctx, db, ID)
	if err != nil {
		logger.Error(fmt.Errorf("error get from DB key %s: %w", ID, err).Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("\n////////////\n%+v\n/////////////\n", alias)

	if alias.IsDeleted() {
		w.WriteHeader(http.StatusGone)
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
