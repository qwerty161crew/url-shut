package main

import (
	"fmt"
	"url-shortener/config"
	pkg "url-shortener/pkg/handlers"

	"github.com/labstack/echo"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации", err)
		return
	}

	port := cfg.Server.Port
	host := cfg.Server.BaseUrl
	if port == "" {
		config.ParseFlags()
		if config.FlagRunAddr != "" {
			port = config.FlagRunAddr
		} else {
			port = ":8080"
		}
	}
	// правильно вас понял?
	if host == "" {
		config.ParseFlags()
		if config.FlagRunAddr != "" {
			port = config.FlagRunAddr
		} else {
			host = "http://127.0.0.1"
		}
	}
	// Правильно ?
	addres := host + port + cfg.Server.AppUrlPrefix
	e := echo.New()
	if config.RedirectHost != "" {
		e.GET(config.RedirectHost+"/:id", pkg.RedirectHandler)
	} else {
		e.GET("/:id", pkg.RedirectHandler)
	}

	e.POST("/", pkg.ShutUrlHandler)
	e.Start(addres)
}
