package middleware

import (
	"fmt"
	"net/http"
	"time"
	"url-shortener/config"
	"url-shortener/internal/service"
	"url-shortener/pkg/logger"

	"github.com/golang-jwt/jwt"
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
func JWTMiddleware(config *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("auth_token")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(config.Security.Salt), nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID, ok := claims["user_id"].(float64)
				if !ok {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
				}
				c.Set("userID", uint(userID))
				return next(c)
			}

			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}
	}
}
