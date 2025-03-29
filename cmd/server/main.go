package main

import (
	"url-shortener/config"
	handlers "url-shortener/internal/handlers"
	md "url-shortener/internal/middlewars"
	"url-shortener/internal/service"
	"url-shortener/pkg/logger"

	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
	"github.com/rs/zerolog/log"
)

// var File string

func main() {
	file := &service.File
	cfg, err := config.LoadConfig()
	if cfg.Server.File != "url.txt" {
		if config.FileUrl != "" {
			*file = config.FileUrl
		}
	} else {
		*file = cfg.Server.File
	}
	if err != nil {
		log.Error().Msg("failed to load config")
		return
	}
	if err := logger.Setup(cfg.Server); err != nil {
		log.Error().Msg("failed to load config")
		return
	}
	err_load := service.LoadData()
	if err_load != nil {
		log.Error().Msg("failed to load data")
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
		if config.FlagRunAddr != "" {
			port = config.FlagRunAddr
		} else {
			host = "http://127.0.0.1"
		}
	}
	addres := host + port
	e := echo.New()
	if config.RedirectHost != "" {
		e.GET(cfg.Server.AppUrlPrefix+config.RedirectHost+"/:id", handlers.RedirectHandler)
	} else {
		e.GET(cfg.Server.AppUrlPrefix+"/:id", handlers.RedirectHandler)
	}
	e.Use(md.LoggerMiddleware)
	e.Use(echoMiddleware.Gzip())
	e.POST("/", handlers.ShutUrlHandler)
	e.POST("/api/shorten", handlers.ShutUrlJsonHandler)
	e.Start(addres)
}
