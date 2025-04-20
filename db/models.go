package db

import (
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetModels() []interface{} {
	return []interface{}{User{}, URL{}}
}

var ErrUsernameExists = errors.New("username already exists")

type User struct {
	gorm.Model
	Username string `gorm:"type:CHAR(16);unique;not null"`
	Password string `gorm:"not null"`
}

type UserRepository interface {
	CreateUser(username string, password string) (User, error)
}
type userRepositoryStruct struct {
	db *gorm.DB
}

func (r *userRepositoryStruct) CreateUser(username string, password string) (User, error) {
	var existingUser User
	result := r.db.Where("username = ?", username).First(&existingUser)

	if result.Error == nil {
		return User{}, ErrUsernameExists
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return User{}, result.Error
	}

	user := User{
		Username: username,
		Password: password,
	}
	if err := r.db.Create(&user).Error; err != nil {
		return User{}, err
	}

	return user, nil
}

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
func NewUserRepository(db *gorm.DB) *userRepositoryStruct {
	return &userRepositoryStruct{db: db}
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
		return nil
	}
	return r.db.Create(&urls).Error
}
func AutoMigrateModels(dbUrl string) error {
	models := GetModels()
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("database connection is nil")
	}
	for _, model := range models {
		if model == nil {
			return fmt.Errorf("one of the models is nil")
		}
	}
	err = db.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("failed to auto migrate models: %v", err)
	}

	return nil
}
