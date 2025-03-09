package pkg

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
	e.POST("/", ShutUrlHandler)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := []byte(test.request)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := ShutUrlHandler(c)
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
	fmt.Println(Urls)
	for key, value := range Urls {
		test := struct {
			name    string
			request string
			want    want
			method  string
		}{
			name:    "Test send request " + value,
			request: "http://127.0.0.1:8080/" + key + "/",
			want:    want{code: 301},
			method:  "GET",
		}
		tests = append(tests, test)
	}
	e := echo.New()
	e.GET("/:id", RedirectHandler)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Println(Urls, test.request)
			body := []byte(test.request)
			req := httptest.NewRequest(http.MethodGet, test.request, bytes.NewBuffer(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := RedirectHandler(c)

			require.NoError(t, err)
			assert.Equal(t, test.want.code, rec.Code)
		})
	}
}
