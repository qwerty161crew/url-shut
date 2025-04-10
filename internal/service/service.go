package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"url-shortener/config"
	"url-shortener/internal/models"
	"url-shortener/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length      = 8
)

type Url struct {
	Id    string `json:"uuid"`
	UrlId string `json:"short_url"`
	Url   string `json:"original_url"`
}
type URLData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var Urls = make(map[string]string)

var File string

type SafeMap struct {
	mu   sync.Mutex
	urls *map[string]string
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		urls: &Urls,
	}
}
func (sm *SafeMap) Set(key, value string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	(*sm.urls)[key] = value
}

func GetUrlInDb(id string, cfg *config.Config) string {
	db, _ := gorm.Open(postgres.Open(cfg.Postgres.GenerateDBurl()), &gorm.Config{})
	urlRepo := models.NewURLRepository(db)
	url, _ := urlRepo.GetBySlug(id)
	return url.Url
}
func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
func generateRandomString() string {
	sb := strings.Builder{}
	sb.Grow(8)

	for i := 0; i < 8; i++ {
		randomIndex := rand.Intn(len(letterBytes))
		sb.WriteByte(letterBytes[randomIndex])
	}

	return sb.String()
}

func SaveUrlInDb(url string, cfg *config.Config) string {
	db, _ := gorm.Open(postgres.Open(cfg.Postgres.GenerateDBurl()), &gorm.Config{})
	urlRepo := models.NewURLRepository(db)
	id := generateRandomString()
	newURL := &models.URL{
		SlugUrl: id,
		Url:     url,
	}
	if err := urlRepo.Create(newURL); err != nil {
		return fmt.Sprintf("Error creating URL: %v\n", err)
	}
	return id
}

func SaveUrlInFile(id string, url string) error {
	file, err := os.OpenFile(File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		logger.Warn("Error open file", err, "file path", File)
		return err
	}
	defer file.Close()
	uuid_id := GenerateUUID()
	var urlstruct Url = Url{Id: uuid_id, UrlId: id, Url: url}
	data, err_serializer := json.Marshal(urlstruct)
	if err_serializer != nil {
		logger.Warn("Error serialize data", err_serializer)
		return err
	}
	data = append(data, '\n')
	_, err = file.Write(data)
	if err != nil {
		logger.Warn("Error writing to file:", err)
		return err
	}
	return nil
}

func LoadData() error {
	file, err := os.OpenFile(File, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		var urlData URLData
		if err := json.Unmarshal(line, &urlData); err != nil {
			return err
		}
		Urls[urlData.ShortURL] = urlData.OriginalURL
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	fmt.Println(Urls)
	return nil
}

func SaveBatchURLS(urls []models.CreateURLSRequest, cfg *config.Config) error {
	gormUrls := make([]models.URL, 0, len(urls))
	for _, url := range urls {
		gormUrls = append(gormUrls, models.URL{
			SlugUrl: url.CorrelationID,
			Url:     url.OriginalURL,
		})
	}
	db, _ := gorm.Open(postgres.Open(cfg.Postgres.GenerateDBurl()), &gorm.Config{})
	urlRepo := models.NewURLRepository(db)
	err := urlRepo.BatchCreate(gormUrls)
	if err != nil {
		return err
	}
	return nil
}
