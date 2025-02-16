package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kamchatkin/practicum-shortener/internal/data"
	"github.com/kamchatkin/practicum-shortener/internal/logs"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"net/http"
	"time"
)

type BatchItemFrom struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchItemTo struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func HandleAPIBatch(w http.ResponseWriter, r *http.Request) {
	var toShort []BatchItemFrom
	logger := logs.NewLogger()
	err := json.NewDecoder(r.Body).Decode(&toShort)
	defer r.Body.Close()
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	var toResp []*BatchItemTo
	var toStore = map[string]string{}

	db, err := storage.NewStorage()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, from := range toShort {
		shortCode, err := getShortCode(ctx, db)
		if err != nil {
			logger.Error(fmt.Errorf("get short code error: %w", err).Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		toResp = append(toResp, &BatchItemTo{
			CorrelationID: from.CorrelationID,
			ShortURL: getShortURL(shortCode, &aliasProps{
				SourceURL: from.OriginalURL,
				HTTPS:     r.TLS != nil,
				Host:      r.Host,
			}),
		})
		toStore[shortCode] = from.OriginalURL
	}

	err = data.SetBatch(ctx, db, toStore)
	if err != nil {
		logger.Error(fmt.Errorf("could not store batch: %w", err).Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toResp)
}
