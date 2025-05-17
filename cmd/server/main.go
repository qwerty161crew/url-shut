package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"url-shortener/config"
	"url-shortener/internal/handlers"
	md "url-shortener/internal/middlewars"
	"url-shortener/internal/repository.go"
	"url-shortener/internal/service"
	"url-shortener/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// var File string
func GetFileStoragePath(cfg *config.Config) string {
	if config.FileUrl != "" {
		return config.FileUrl
	}
	return cfg.Server.File
}

func main() {
	cfg, err := config.LoadConfig()
	dbConnect, err := gorm.Open(postgres.Open(cfg.Postgres.GenerateDBurl()), &gorm.Config{})
	userRepository := repository.NewUserRepository(dbConnect)
	urlRepository := repository.NewURLRepository(dbConnect)
	urlService := service.GetUrlService(userRepository, urlRepository)
	if err != nil {
		log.Error().Msg("failed to open connect db")
		return 
	}
	if err != nil {
		log.Error().Msg("failed to load config")
		return
	}

	if err := repository.AutoMigrateModels(cfg.Postgres.GenerateDBurl()); err != nil {
		log.Error().Msg(err.Error())
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
		if config.FlagRunAddr != "" {
			port = config.FlagRunAddr
		} else {
			host = "http://127.0.0.1"
		}
	}
	addres := host + port
	e := echo.New()
	urlHandler := handlers.NewURLHandler(cfg,  urlService)
	if config.RedirectHost != "" {
		e.GET(cfg.Server.AppUrlPrefix+config.RedirectHost+"/:id", urlHandler.RedirectHandler)
	} else {
		e.GET(cfg.Server.AppUrlPrefix+"/:id", urlHandler.RedirectHandler)
	}
	e.Use(md.LoggerMiddleware)
	e.Use(echoMiddleware.Gzip())
	e.Use(md.JWTMiddleware(cfg))
	e.POST("/", urlHandler.ShutUrlHandler)
	e.POST("/api/shorten", urlHandler.ShutUrlJsonHandler)
	e.POST("/api/shorten/batch", urlHandler.BatchURLHandler)
	e.POST("/api/auth", urlHandler.RegistrationHandler)
	e.GET("/ping", func(ctx echo.Context) error {
		connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Db)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			logger.Error("connect db error", err)
			return ctx.String(http.StatusInternalServerError, "connect db error")
		}
		if err := db.Ping(); err != nil {
			logger.Error("connect db error", err)
			return ctx.String(http.StatusInternalServerError, "connect db error")
		}

		defer db.Close()
		return ctx.String(http.StatusOK, "Success")
	})
	e.Start(addres)
}
