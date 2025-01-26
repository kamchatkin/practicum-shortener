package app

import (
	"encoding/json"
	"net/http"
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
	err := json.NewDecoder(r.Body).Decode(&toShort)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validate.Struct(&toShort)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := makeAlias(&aliasProps{
		SourceURL: toShort.URL,
		HTTPS:     r.TLS != nil,
		Host:      r.Host,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(APIResponse{Result: shortURL})
}
