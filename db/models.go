package db

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	SlugUrl string `gorm:"type:CHAR(8);unique;not null"`
	Url     string `gorm:"type:TEXT;not null"`
}
type URLRepository interface {
	GetBySlug(slug string) (URL, error)
	Create(url *URL) error
	BatchCreate(urls []URL) error
	CreateOrGet(url *URL) (*URL, error)
}

type urlRepository struct {
	db *gorm.DB
}

func (r *urlRepository) CreateOrGet(url *URL) (*URL, error) {
	if url == nil {
		return nil, errors.New("URL cannot be nil")
	}

	// Пытаемся вставить новую запись
	result := r.db.Create(url)
	if result.Error == nil {
		return url, nil // Успешная вставка
	}

	// Проверяем, является ли ошибка нарушением уникальности
	var pgErr *pgconn.PgError
	if errors.As(result.Error, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		// Если URL уже существует, находим существующую запись
		var existingURL URL
		if err := r.db.Where("url = ?", url.Url).First(&existingURL).Error; err != nil {
			return nil, err
		}
		return &existingURL, gorm.ErrDuplicatedKey
	}

	return nil, result.Error
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
