package config

import (
	"flag"
	"fmt"
	"net/url"
)

type ConfigType struct {
	// Addr адрес для запуска этого приложения
	Addr string

	// ShortHost подменный хост в сокращенном УРЛ
	ShortHost string

	ShortHostURL *url.URL
}

var Config = ConfigType{}

func Parse() {
	if Config != (ConfigType{}) {
		return
	}

	flag.StringVar(&Config.Addr, "a", ":8080", "Адрес запуска сервера. HOST:PORT")
	flag.StringVar(&Config.ShortHost, "b", "", "Подменный УРЛ для сокращенного УРЛ. HOST:PORT")
	flag.Parse()

	if Config.ShortHost != "" {
		parsedURL, err := url.Parse(Config.ShortHost)
		if err != nil {
			fmt.Println("Не удалось разобрать URL для подстановки в сокращенную ссылку")
			panic(err)
		}

		Config.ShortHostURL = parsedURL
	}
}
