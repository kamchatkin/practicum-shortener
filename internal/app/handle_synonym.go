package app

import (
	"fmt"
	"io"
	"net/http"
)

const maxIterate = 3

// SynonymHandler Создание сокращения
func SynonymHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sourceURL, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validate.Var(string(sourceURL), "url")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var shortURL string
	shortURL, err = makeAlias(&aliasProps{
		SourceURL: string(sourceURL),
		HTTPS:     r.TLS != nil,
		Host:      r.Host,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprint(w, shortURL)
}
