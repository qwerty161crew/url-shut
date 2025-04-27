package service

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
	"url-shortener/config"
	"url-shortener/db"
	"url-shortener/internal/models"
	"url-shortener/internal/repository.go"

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

type UrlService struct {
	db  *gorm.DB
}

func GetUrlService(db gorm.DB) *UrlService {
	return &UrlService{
		db: &db,
	}
}
func (s *UrlService) CreateUser(user models.RegistrationRequest, cfg *config.Config) (string, error) {
	passwordAndsalt := user.Password + cfg.Security.Salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordAndsalt), bcrypt.DefaultCost)
	if err != nil {

		return "", fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}

	userRepository := repository.NewUserRepository(s.db)
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

func (s *UrlService) GetUrlsInDb(userID uint, cfg *config.Config) []models.ListURLSResponse {
	var responses []models.ListURLSResponse
	urlRepo := repository.NewURLRepository(s.db)
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

func (s *UrlService) GetUrlInDb(id string, cfg *config.Config) string {
	urlRepo := repository.NewURLRepository(s.db)
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

func (s *UrlService) SaveUrlInDb(url string, cfg *config.Config, userId uint) (string, error) {

	urlRepo := repository.NewURLRepository(s.db)
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


func (s *UrlService) SaveBatchURLS(urls []models.CreateURLSRequest, cfg *config.Config) ([]models.CreateURLSResponse, error) {
	gormUrls := make([]db.URL, 0, len(urls))
	var responses []models.CreateURLSResponse
	for _, url := range urls {
		gormUrls = append(gormUrls, db.URL{
			Slug: url.CorrelationID,
			Url:     url.OriginalURL,
		})
	}
	urlRepo := repository.NewURLRepository(s.db)
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
