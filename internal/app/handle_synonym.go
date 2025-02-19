package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/kamchatkin/practicum-shortener/internal/logs"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"io"
	"net/http"
	"time"
)

const maxIterate = 3

// SynonymHandler Создание сокращения
func SynonymHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	logger := logs.NewLogger()

	sourceURL, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validate.Var(string(sourceURL), "url")
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := storage.NewStorage()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	alProps := &aliasProps{
		SourceURL: string(sourceURL),
		HTTPS:     r.TLS != nil,
		Host:      r.Host,
	}

	var shortURL string
	shortURL, err = makeAlias(ctx, db, alProps)
	if err != nil {
		if errors.Is(err, ErrUniq) {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(shortURL))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprint(w, shortURL)
}
