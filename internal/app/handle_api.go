package app

import (
	"context"
	"encoding/json"
	"github.com/kamchatkin/practicum-shortener/internal/logs"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"net/http"
	"time"
)

// APIResponse Структура ответа
type APIResponse struct {
	Result string `json:"result" validate:"required,uri"`
}

type APIRequest struct {
	URL string `json:"url" validate:"required,uri"`
}

// HandleAPI Сокращение ссылки по api
func HandleAPI(w http.ResponseWriter, r *http.Request) {
	toShort := APIRequest{}
	logger := logs.NewLogger()
	err := json.NewDecoder(r.Body).Decode(&toShort)
	defer r.Body.Close()
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validate.Struct(&toShort)
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
		SourceURL: toShort.URL,
		HTTPS:     r.TLS != nil,
		Host:      r.Host,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	shortURL, err := makeAlias(ctx, db, alProps)
	if err != nil {
		switch err {
		case ErrUniq:
			w.WriteHeader(http.StatusConflict)
		default:
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	_ = json.NewEncoder(w).Encode(APIResponse{Result: shortURL})
}
