package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortener/internal/handlers"
	"url-shortener/internal/models"
	"url-shortener/internal/service"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShutUrlHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name    string
		request string
		want    want
		method  string
	}{
		{
			name:    "positive test #1",
			request: "https://practicum.yandex.ru",
			want: want{
				code:        201,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			name:    "Test bad request",
			request: "Invalid url",
			want: want{
				code:        400,
				response:    "",
				contentType: "text/plain",
			},
		},
	}

	e := echo.New()
	e.POST("/", handlers.ShutUrlHandler)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := []byte(test.request)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handlers.ShutUrlHandler(c)
			require.NoError(t, err)
			assert.Equal(t, test.want.code, rec.Code)
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		name    string
		request string
		want    want
		method  string
	}{}

	for key, value := range service.Urls {
		test := struct {
			name    string
			request string
			want    want
			method  string
		}{
			name:    "Test send request " + value,
			request: "/" + key, // Используем только путь, а не полный URL
			want:    want{code: 301},
			method:  "GET",
		}
		tests = append(tests, test)
	}

	e := echo.New()
	e.GET("/:id", handlers.RedirectHandler)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Println(service.Urls, test.request)

			req := httptest.NewRequest(http.MethodGet, test.request, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(strings.TrimPrefix(test.request, "/"))

			err := handlers.RedirectHandler(c)

			require.NoError(t, err)
			assert.Equal(t, test.want.code, rec.Code)
		})
	}
}

func TestShutUrlJsonHandler(t *testing.T) {
	e := echo.New()
	t.Run("Valid POST request", func(t *testing.T) {
		requestBody := models.RequestCreateUrl{Url: "https://example.com"}
		jsonBody, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(jsonBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, handlers.ShutUrlJsonHandler(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)

			var response models.ResponseCreateUrl
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response.Result, "http://")
		}
	})
}
