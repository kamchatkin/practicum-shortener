package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/kamchatkin/practicum-shortener/internal/app"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	// Сокращение
	r.Post("/", app.SynonymHandler)

	// Переадресация
	r.Get("/{id}", app.RedirectHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
