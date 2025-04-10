package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"url-shortener/config"
	"url-shortener/internal/models"
	"url-shortener/internal/service"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type UrlHandlers struct {
	db gorm.DB
}

func NewUrlHandlers(db gorm.DB) UrlHandlers {
	return UrlHandlers{db: db}
}

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

type URLHandler struct {
	config *config.Config
}

func NewURLHandler(cfg *config.Config) *URLHandler {
	return &URLHandler{
		config: cfg,
	}
}

func (h *URLHandler) ShutUrlJsonHandler(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.String(http.StatusMethodNotAllowed, "Only POST requests are allowed!")
	}

	var request models.RequestCreateUrl
	bodyBytes, _ := io.ReadAll(c.Request().Body)
	err := json.Unmarshal(bodyBytes, &request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	id := service.SaveUrlInDb(request.Url, h.config)
	link := fmt.Sprintf("%s://%s/%s", c.Scheme(), c.Request().Host, id)
	response := models.ResponseCreateUrl{Result: link}
	return c.JSON(http.StatusCreated, response)
}

func (h *URLHandler) ShutUrlHandler(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.String(http.StatusMethodNotAllowed, "Only POST requests are allowed!")
	}

	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read request body")
	}
	defer c.Request().Body.Close()

	body := string(bodyBytes)
	if !isValidURL(body) {
		return c.String(http.StatusBadRequest, "Invalid url")
	}

	id := service.SaveUrlInDb(body, h.config)
	link := fmt.Sprintf("%s://%s/%s", c.Scheme(), c.Request().Host, id)
	return c.String(http.StatusCreated, link)
}

func (h *URLHandler) RedirectHandler(c echo.Context) error {
	originalURL := service.GetUrlInDb(c.Param("id"), h.config)
	originalURL = strings.TrimSpace(originalURL)
	originalURL = strings.Trim(originalURL, `"`)
	if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
		originalURL = `http://` + originalURL
	}

	return c.Redirect(http.StatusMovedPermanently, originalURL)
}

func (h *URLHandler) BatchURLHandler(c echo.Context) error {
	// var typereq models.CreateURLSRequest
	var requests []models.CreateURLSRequest
	if err := c.Bind(&requests); err != nil {
		fmt.Println("24432432")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	// if err := c.Validate(typereq); err != nil {
	// 	fmt.Println("555")
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	// }
	response, err := service.SaveBatchURLS(requests, h.config)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	return c.JSON(http.StatusBadRequest, response)
}
