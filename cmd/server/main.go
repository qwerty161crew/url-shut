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
	if port == "" {
		config.ParseFlags()
		if config.FlagRunAddr != "" {
			port = config.FlagRunAddr
		} else {
			port = ":8080"
		}
	}

	e := echo.New()
	e.POST("/", pkg.ShutUrlHandler)
	e.GET("/:id", pkg.RedirectHandler)
	e.Start(port)
}
