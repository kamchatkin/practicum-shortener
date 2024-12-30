package app

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

var db = map[string]string{
	"qwerty": "http://localhost:8080/?qwerty",
}

// RedirectHandler Поиск сокращения и переадресация
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "id")

	var toURL string
	var ok bool
	if toURL, ok = db[alias]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", toURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
