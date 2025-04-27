package service

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
	"url-shortener/config"
	"url-shortener/db"
	"url-shortener/internal/models"
	"url-shortener/internal/repository.go"
	"url-shortener/pkg/logger"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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

func CreateUser(user models.RegistrationRequest, cfg *config.Config, dbConnect *gorm.DB) (string, error) {

	passwordAndsalt := user.Password + cfg.Security.Salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordAndsalt), bcrypt.DefaultCost)
	if err != nil {

		return "", fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}

	userRepository := repository.NewUserRepository(dbConnect)
	userModel, err := userRepository.CreateUser(user.Username, string(hashedPassword))
	fmt.Println(userModel)
	if err != nil {
		return "", err
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userModel.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, _ := token.SignedString([]byte(cfg.Security.Salt))
	return tokenString, nil
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

func GetUrlsInDb(userID uint, cfg *config.Config, dbConnect *gorm.DB) []models.ListURLSResponse {
	var responses []models.ListURLSResponse
	urlRepo := repository.NewURLRepository(dbConnect)
	urls, _ := urlRepo.GetListUrls(userID)
	for _, url := range urls {
		response := models.ListURLSResponse{
			ShortUrl:    fmt.Sprintf("http://%s/%s:%s", cfg.Server.BaseUrl, cfg.Server.Port, url.Slug),
			OriginalUrl: url.Url,
		}
		responses = append(responses, response)
	}
	return responses
}

func GetUrlInDb(id string, cfg *config.Config, dbConnect *gorm.DB) string {
	urlRepo := repository.NewURLRepository(dbConnect)
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

func SaveUrlInDb(url string, cfg *config.Config, userId uint, dbConnect *gorm.DB) (string, error) {

	urlRepo := repository.NewURLRepository(dbConnect)
	id := generateRandomString()
	newURL := &db.URL{
		Slug: id,
		Url:     url,
		UserID:  userId,
	}

	createdURL, err := urlRepo.CreateIfNotExists(newURL)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return createdURL.Slug, err
		}
		return "", err
	}

	return createdURL.Slug, nil
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

func SaveBatchURLS(urls []models.CreateURLSRequest, cfg *config.Config, dbConnect *gorm.DB) ([]models.CreateURLSResponse, error) {
	gormUrls := make([]db.URL, 0, len(urls))
	var responses []models.CreateURLSResponse
	for _, url := range urls {
		gormUrls = append(gormUrls, db.URL{
			Slug: url.CorrelationID,
			Url:     url.OriginalURL,
		})
	}
	urlRepo := repository.NewURLRepository(dbConnect)
	err := urlRepo.BatchCreate(gormUrls)
	if err != nil {
		return nil, err
	}

	for _, req := range urls {
		response := models.CreateURLSResponse{
			CorrelationID: req.CorrelationID,
			OriginalURL:   fmt.Sprintf("http://%s/%s:%s", cfg.Server.BaseUrl, cfg.Server.Port, req.CorrelationID),
		}
		responses = append(responses, response)
	}
	return responses, nil
}
