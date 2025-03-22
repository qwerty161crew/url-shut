package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"url-shortener/internal/models"
	service "url-shortener/internal/service"

	"github.com/labstack/echo"
)

func isValidURL(urlString string) bool {
	u, err := url.ParseRequestURI(urlString)
	fmt.Println(u, err)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	return true
}
func ShutUrlJsonHandler(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.String(http.StatusMethodNotAllowed, "Only POST requests are allowed!")
	}
	var request models.RequestCreateUrl
	bodyBytes, _ := io.ReadAll(c.Request().Body)
	err := json.Unmarshal(bodyBytes, &request)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	id := service.SaveUrl(request.Url)
	link := fmt.Sprintf("%s://%s/%s", c.Scheme(), c.Request().Host, id)
	response := models.ResponseCreateUrl{Result: link}
	return c.JSON(http.StatusCreated, response)

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
	id := service.SaveUrl(body)
	link := fmt.Sprintf("%s://%s/%s", c.Scheme(), c.Request().Host, id)
	return c.String(http.StatusCreated, link)
}

func RedirectHandler(c echo.Context) error {
	originalURL, exists := service.Urls[c.Param("id")]
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
