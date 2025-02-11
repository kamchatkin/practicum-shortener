package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kamchatkin/practicum-shortener/config"
	"net/http"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.Config()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseDsn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer dbpool.Close()

	var one int
	err = dbpool.QueryRow(context.Background(), "select 1").Scan(&one)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
