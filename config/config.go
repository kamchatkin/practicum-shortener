package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v11"
	"net/url"
)

type ConfigType struct {
	// Addr адрес для запуска этого приложения
	Addr string `env:"SERVER_ADDRESS"`

	// ShortHost подменный хост в сокращенном УРЛ
	ShortHost string `env:"BASE_URL"`

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

	//if Config.ShortHost != "" {
	//	confShortHost(Config.ShortHost)
	//}

	parseEnv()
}

func parseEnv() {
	//if _addr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
	//	Config.Addr = _addr
	//}
	//
	//if _shortHost, ok := os.LookupEnv("BASE_URL"); ok {
	//	confShortHost(_shortHost)
	//}

	err := env.Parse(&Config)
	if err != nil {
		panic(err)
	}

	if Config.ShortHost != "" {
		confShortHost(Config.ShortHost)
	}
}

func confShortHost(shortHost string) {
	if shortHost == "" {
		return
	}

	parsedURL, err := url.Parse(shortHost)
	if err != nil || parsedURL.Host == "" {
		fmt.Println(fmt.Sprintf("Не удалось разобрать URL для подстановки в сокращенную ссылку. Получено: %s", shortHost))
		panic(err)
	}

	Config.ShortHost = shortHost
	Config.ShortHostURL = parsedURL
}
