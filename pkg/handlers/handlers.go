package pkg

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length      = 8
)

var Urls = make(map[string]string)

func ShutUrlHandler(c echo.Context) error {

	if c.Request().Method != http.MethodPost {
		return c.String(http.StatusMethodNotAllowed, "Only POST requests are allowed!")
	}

	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read request body")
	}
	defer c.Request().Body.Close()
	body := string(bodyBytes)
	id := generateRandomString(length)
	Urls[id] = body
	fmt.Println(Urls)

	link := fmt.Sprintf("http://127.0.0.1/%s", id)
	return c.String(http.StatusCreated, link)
}
func generateRandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)

	for i := 0; i < n; i++ {
		randomIndex := rand.Intn(len(letterBytes))
		sb.WriteByte(letterBytes[randomIndex])
	}

	return sb.String()
}
func RedirectHandler(c echo.Context) error {
	originalURL, exists := Urls[c.Param("id")]
	fmt.Println(c.Param("id"), "342324234")

	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "url not found")
	}
	originalURL = strings.TrimSpace(originalURL)
	originalURL = strings.Trim(originalURL, `"`)
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = `http://` + originalURL
	}

	return c.Redirect(http.StatusMovedPermanently, originalURL)
}
