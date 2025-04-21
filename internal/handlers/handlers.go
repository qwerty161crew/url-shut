package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"url-shortener/config"
	"url-shortener/db"
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
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "user ID not found in context"})
	}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	id, err := service.SaveUrlInDb(request.Url, h.config, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			link := fmt.Sprintf("%s://%s/%s", c.Scheme(), c.Request().Host, id)
			return c.JSON(http.StatusConflict, models.ResponseCreateUrl{Result: link})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	link := fmt.Sprintf("%s://%s/%s", c.Scheme(), c.Request().Host, id)
	return c.JSON(http.StatusCreated, models.ResponseCreateUrl{Result: link})
}
func (h *URLHandler) ShutUrlHandler(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.String(http.StatusMethodNotAllowed, "Only POST requests are allowed!")
	}
	userID, _ := c.Get("userID").(uint)
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read request body")
	}
	defer c.Request().Body.Close()

	body := string(bodyBytes)
	if !isValidURL(body) {
		return c.String(http.StatusBadRequest, "Invalid url")
	}

	id, _ := service.SaveUrlInDb(body, h.config, userID)

	link := fmt.Sprintf("%s://%s/%s", c.Scheme(), c.Request().Host, id)
	// GET /api/user/urls
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
func (h *URLHandler) RegistrationHandler(c echo.Context) error {
	var request models.RegistrationRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	fmt.Println(request)
	token, err := service.CreateUser(request, h.config)
	if err != nil {
		if errors.Is(err, db.ErrUsernameExists) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "username already exists"})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Path = "/"

	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, map[string]string{"message": "registration was successful"})
}

func (h *URLHandler) GetUserUrls(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "user ID not found in context"})
	}
	urls := service.GetUrlsInDb(userID, h.cfg)
	return c.JSON(http.StatusOK, urls)
}
