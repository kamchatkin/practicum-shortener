package main

import (
	"fmt"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/router"
	"net/http"
)

func main() {
	cfg, err := config.Config()
	if err != nil {
		fmt.Printf("Ошибка подготовки конфигурации приложения. Надо ли в этом случае давать запускать приложение?\n%s", err)
		panic(err)
	}

	if err := http.ListenAndServe(cfg.Addr, router.Router()); err != nil {
		panic(err)
	}
}
