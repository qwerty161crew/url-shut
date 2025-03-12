package main

import (
	"url-shortener/config"
	"url-shortener/internal/handlers"
	internal "url-shortener/internal/handlers"
	"url-shortener/pkg/logger"

	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error().Msg("failed to load config")
		return
	}
	if err := logger.Setup(cfg.Server); err != nil {
		log.Error().Msg("failed to load config")
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
	if host == "" {
		config.ParseFlags()
		if config.FlagRunAddr != "" {
			port = config.FlagRunAddr
		} else {
			host = "http://127.0.0.1"
		}
	}
	addres := host + port + cfg.Server.AppUrlPrefix
	e := echo.New()
	if config.RedirectHost != "" {
		e.GET(config.RedirectHost+"/:id", internal.RedirectHandler)
	} else {
		e.GET("/:id", handlers.RedirectHandler)
	}

	e.POST("/", handlers.ShutUrlHandler)
	e.Start(addres)
}
