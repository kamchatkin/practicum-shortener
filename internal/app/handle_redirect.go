package app

import (
	"net/http"
)

var db = map[string]string{
	"qwerty": "http://localhost:8080/?qwerty",
}

// RedirectHandler Поиск сокращения и переадресация
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	alias := r.PathValue("id")
	var toURL string
	var ok bool
	if toURL, ok = db[alias]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", toURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
