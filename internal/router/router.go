package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/kamchatkin/practicum-shortener/internal/app"
	"github.com/kamchatkin/practicum-shortener/internal/router/middlewares/gzip"
	"github.com/kamchatkin/practicum-shortener/internal/router/middlewares/log"
	"net/http"
)

// Router Маршруты приложения
func Router() *chi.Mux {
	r := chi.NewRouter()

	//r.Use(middleware.Compress(5, ""))

	// Сокращение
	r.Post("/", handleWrapper(app.SynonymHandler))

	// Переадресация
	r.Get("/{id}", handleWrapper(app.RedirectHandler))

	// api, iter7
	r.Post("/api/shorten", handleWrapper(app.HandleAPI))

	return r
}

// handleWrapper обертка хэндлеров в мидлвари
func handleWrapper(next http.HandlerFunc) http.HandlerFunc {
	mwsList := []func(next http.HandlerFunc) http.HandlerFunc{
		log.WithLogging,
		gzip.WithGzipped,
	}

	for _, middleware := range mwsList {
		next = middleware(next)
	}

	return next
}
