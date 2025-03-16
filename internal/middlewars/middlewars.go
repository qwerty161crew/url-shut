package middleware

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"
	"url-shortener/pkg/logger"

	"github.com/labstack/echo"
)

type LogWrite struct {
	uri    string
	method string
	time   string
}

func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := GenerateUUID()
		start := time.Now()
		logger.Info(
			"Request started",
			"request_id", requestID,
			"method", c.Request().Method,
			"uri", c.Request().RequestURI,
		)

		err := next(c)

		duration := time.Since(start)
		if err != nil {
			logger.Warn(
				"Request failed",
				"request_id", requestID,
				"method", c.Request().Method,
				"uri", c.Request().RequestURI,
				"status", c.Response().Status,
				"duration", duration,
				"error", err,
			)
		} else {
			logger.Info(
				"Request completed",
				"request_id", requestID,
				"method", c.Request().Method,
				"uri", c.Request().RequestURI,
				"status", c.Response().Status,
				"size", c.Response().Size,
				"duration", duration,
			)
		}

		return err
	}
}
