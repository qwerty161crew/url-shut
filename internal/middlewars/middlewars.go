package middleware

import (
	"time"
	"url-shortener/internal/service"
	"url-shortener/pkg/logger"

	"github.com/labstack/echo"
)

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestID := service.GenerateUUID()
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
