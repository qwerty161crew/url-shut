package service

import (
	"math/rand"
	"strings"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length      = 8
)

var Urls = make(map[string]string)

// generateRandomString генерирует случайную строку заданной длины.
func generateRandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)

	for i := 0; i < n; i++ {
		randomIndex := rand.Intn(len(letterBytes))
		sb.WriteByte(letterBytes[randomIndex])
	}

	return sb.String()
}

// SaveUrl сохраняет URL и возвращает сгенерированный идентификатор.
func SaveUrl(url string) string {
	id := generateRandomString(length)
	Urls[id] = url
	return id
}
