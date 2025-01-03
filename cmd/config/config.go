package config

import (
	"flag"
)

type ConfigType struct {
	// Addr адрес для запуска этого приложения
	Addr string

	// ShortHost подменный хост в сокращенном УРЛ
	ShortHost string
}

var Config = ConfigType{}

func Parse() {
	if Config != (ConfigType{}) {
		return
	}

	flag.StringVar(&Config.Addr, "a", ":8080", "Адрес запуска сервера. HOST:PORT")
	flag.StringVar(&Config.ShortHost, "b", "", "Подменный УРЛ для сокращенного УРЛ. HOST:PORT")
	flag.Parse()
}
