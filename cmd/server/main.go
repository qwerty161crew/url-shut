package main

import (
	"fmt"
	"url-shortener/config"
	pkg "url-shortener/pkg/handlers"

	"github.com/labstack/echo"
)

func main() {
	cfg, err := config.LoadConfig()
	fmt.Println(cfg)
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации", err)
	}
	e := echo.New()
	e.POST("/", pkg.ShutUrlHandler)
	e.GET("/:id", pkg.RedirectHandler)
	e.Start(cfg.Server.Port)
}
