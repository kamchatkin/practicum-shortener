package main

import (
	"github.com/kamchatkin/practicum-shortener/internal/app"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Переадресация
	mux.HandleFunc("GET /{id}", app.RedirectHandler)

	// Сокращение
	mux.HandleFunc("POST /", app.SynonymHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
