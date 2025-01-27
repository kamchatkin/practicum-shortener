package app

import (
	"bufio"
	"encoding/json"
	"github.com/kamchatkin/practicum-shortener/config"
	"os"
)

var db = map[string]string{
	"qwerty": "http://localhost:8080/?qwerty",
}

type dbRecord struct {
	Alias  string `json:"alias"`
	Source string `json:"source"`
}

// SaveDB Дамп БД на диск
func SaveDB() {

	cfg, err := config.Config()
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(cfg.DumpPath(), []byte{}, 0666)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(cfg.DumpPath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for k, v := range db {
		rec := dbRecord{
			Alias:  k,
			Source: v,
		}
		_ = json.NewEncoder(writer).Encode(&rec)
	}
}

// LoadDB Чтение дампа БД
func LoadDB() {
	cfg, err := config.Config()
	if err != nil {
		panic(err)
	}

	file, err := os.Open(cfg.DumpPath())
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rec := dbRecord{}
		_ = json.Unmarshal([]byte(scanner.Text()), &rec)
		db[rec.Alias] = rec.Source
	}
}
