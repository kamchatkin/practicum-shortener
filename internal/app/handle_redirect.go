package app

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// RedirectHandler Поиск сокращения и переадресация
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")

	var toURL string
	var ok bool
	if toURL, ok = db[alias]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Location", toURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
