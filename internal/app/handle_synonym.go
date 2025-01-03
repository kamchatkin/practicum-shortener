package app

import (
	"fmt"
	"github.com/kamchatkin/practicum-shortener/config"
	"io"
	"net/http"
	"net/url"
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

	if _, err := url.ParseRequestURI(string(sourceURL)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var alias string
	for i := range maxIterate {
		alias = shortness()

		if _, ok := db[alias]; !ok {
			break
		}

		i++
		if i == maxIterate {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	db[alias] = string(sourceURL)

	w.WriteHeader(http.StatusCreated)

	if config.Config.ShortHost != "" {
		//if strings.HasPrefix(config.Config.ShortHost, "http") {
		fmt.Fprintf(w, "%s://%s/%s", config.Config.ShortHostURL.Scheme, config.Config.ShortHostURL.Host, alias)
		//}

		return
	}

	proto := "http"
	if r.TLS != nil {
		proto = "https"
	}
	_, _ = fmt.Fprintf(w, "%s://%s/%s", proto, r.Host, alias)
}
