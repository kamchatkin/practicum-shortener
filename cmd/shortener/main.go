package main

import (
	"fmt"
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/app"
	"github.com/kamchatkin/practicum-shortener/internal/router"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewDevelopment()
}

func main() {
	ch := make(chan os.Signal, 3)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		signalType := <-ch
		signal.Stop(ch)
		logger.Info("Exit command received. Exiting...")

		// this is a good place to flush everything to disk
		// before terminating.
		app.SaveDB()
		logger.Info("Signal type : ", zap.String("signal", signalType.String()))
		logger.Sync()

		os.Exit(0)
	}()

	cfg, err := config.Config()
	if err != nil {
		fmt.Printf("Ошибка подготовки конфигурации приложения. Надо ли в этом случае давать запускать приложение?\n%s", err)
		panic(err)
	}

	app.LoadDB()
	defer app.SaveDB()
	if err := http.ListenAndServe(cfg.Addr, router.Router()); err != nil {
		panic(err)
	}

}
