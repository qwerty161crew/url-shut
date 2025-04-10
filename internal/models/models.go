package models

import (
	"errors"

	"gorm.io/gorm"
)

const (
	TypeSimpleUtterance = "SimpleUtterance"
)

type CreateURLSRequest struct {
	CorrelationID string `json:"correlation_id" validate:"required"`
	OriginalURL   string `json:"original_url" validate:"required,url"`
}

type CreateURLSResponse struct {
	CorrelationID string `json:"correlation_id" validate:"required"`
	OriginalURL   string `json:"short_url" validate:"required,url"`
}

type RequestCreateUrl struct {
	Url string `json:"url"`
}

type ResponseCreateUrl struct {
	Result string `json:"result"`
}

type URLRepository interface {
	GetBySlug(slug string) (URL, error)
	Create(url *URL) error
	BatchCreate(urls []URL) error
}

type urlRepository struct {
	db *gorm.DB
}

type URL struct {
	gorm.Model
	SlugUrl string `gorm:"type:CHAR(8);unique;not null"`
	Url     string `gorm:"type:TEXT;not null"`
}

func NewURLRepository(db *gorm.DB) URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) GetBySlug(slug string) (URL, error) {
	var url URL
	err := r.db.Where("slug_url = ?", slug).First(&url).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return URL{}, errors.New("URL not found")
		}
		return URL{}, err
	}
	return url, nil
}

func (r *urlRepository) Create(url *URL) error {
	if url == nil {
		return errors.New("URL cannot be nil")
	}
	return r.db.Create(url).Error
}

func (r *urlRepository) BatchCreate(urls []URL) error {
	if len(urls) == 0 {
		return nil // Нет данных для вставки
	}
	return r.db.Create(&urls).Error
}
