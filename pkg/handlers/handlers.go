package pkg

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	pkg "url-shortener/pkg/service"

	"github.com/labstack/echo"
)

func isValidURL(urlString string) bool {
	u, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	if !strings.Contains(u.Host, ".") {
		return false
	}

	return true
}
func ShutUrlHandler(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.String(http.StatusMethodNotAllowed, "Only POST requests are allowed!")
	}

	bodyBytes, err := io.ReadAll(c.Request().Body)
	body := string(bodyBytes)
	defer c.Request().Body.Close()
	isUrl := isValidURL(body)
	if isUrl == false {
		return c.String(http.StatusBadRequest, "Invalid url")
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read request body")
	}
	id := pkg.SaveUrl(body)
	link := fmt.Sprintf("%s://%s/%s", c.Scheme(), c.Request().Host, id)
	return c.String(http.StatusCreated, link)
}

func RedirectHandler(c echo.Context) error {
	originalURL, exists := pkg.Urls[c.Param("id")]
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
