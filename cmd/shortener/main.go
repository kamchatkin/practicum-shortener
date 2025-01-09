package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/app"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	// Сокращение
	r.Post("/", app.SynonymHandler)

	// Переадресация
	r.Get("/{id}", app.RedirectHandler)

	cfg, err := config.Config()
	if err != nil {
		fmt.Printf("Ошибка подготовки конфигурации приложения. Надо ли в этом случае давать запускать приложение?\n%s", err)
		panic(err)
	}

	if err := http.ListenAndServe(cfg.Addr, r); err != nil {
		panic(err)
	}
}
