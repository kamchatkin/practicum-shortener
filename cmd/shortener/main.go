package main

import (
	"github.com/kamchatkin/practicum-shortener/config"
	"github.com/kamchatkin/practicum-shortener/internal/logs"
	"github.com/kamchatkin/practicum-shortener/internal/router"
	"github.com/kamchatkin/practicum-shortener/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := logs.NewLogger()

	ch := make(chan os.Signal, 3)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		signalType := <-ch
		signal.Stop(ch)
		logger.Info("Exit command received. Exiting...")

		// this is a good place to flush everything to disk
		// before terminating.
		logger.Info("Signal type : ", zap.String("signal", signalType.String()))
		logger.Sync()
		storage.DB.Close()
		os.Exit(0)
	}()

	// Конфигурация
	cfg, err := config.Config()
	if err != nil {
		logger.Error("Ошибка подготовки конфигурации приложения. Надо ли в этом случае давать запускать приложение?\n",
			zap.Error(err))
		os.Exit(1)
	}
	// Далее ошибку при получении конфигурации можно игнорировать
	// @todo по хорошему рефакторинг, отдельный метод подготовки настроек и потом только получать объект

	// Подготовка хранилища
	storage.InitStorage()
	if ok, err1 := storage.DB.Open(); !ok {
		if err1 != nil {
			logger.Error(err1.Error())
		}

		os.Exit(1)
	}

	if err := http.ListenAndServe(cfg.Addr, router.Router()); err != nil {
		panic(err)
	}
}
