package app

import (
	"context"
	"github.com/kamchatkin/practicum-shortener/internal/auth"
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
		logger.Info("1")
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validate.Var(string(sourceURL), "url")
	if err != nil {
		logger.Info("2")
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := storage.NewStorage()
	if err != nil {
		logger.Error(err.Error())
		logger.Info("3")
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

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	var shortURL string
	shortURL, err = makeAlias(ctx, db, alProps, auth.GetUserIDFromCookie(r))
	if err != nil {
		switch err {
		case ErrUniq:
			logger.Info("4")
			w.WriteHeader(http.StatusConflict)
		default:
			logger.Info("5")
			logger.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		logger.Info("6")
		w.WriteHeader(http.StatusCreated)
	}

	w.Write([]byte(shortURL))
}
