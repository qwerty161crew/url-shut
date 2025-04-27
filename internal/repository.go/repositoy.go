package repository

import (
	"errors"
	"fmt"
	"url-shortener/db"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(username string, password string) (db.User, error)
}
type userRepositoryStruct struct {
	db *gorm.DB
}

func (r *userRepositoryStruct) CreateUser(username string, password string) (db.User, error) {
	var existingUser db.User
	result := r.db.Where("username = ?", username).First(&existingUser)

	if result.Error == nil {
		return db.User{}, db.ErrUsernameExists
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return db.User{}, result.Error
	}

	user := db.User{
		Username: username,
		Password: password,
	}
	if err := r.db.Create(&user).Error; err != nil {
		return db.User{}, err
	}

	return user, nil
}

type URLRepository interface {
	GetBySlug(slug string) (db.URL, error)
	Create(url *db.URL) error
	BatchCreate(urls []db.URL) error
	CreateIfNotExists(url *db.URL) (*db.URL, error)
	GetListUrls(userId uint) ([]*db.URL, error)
}

type urlRepository struct {
	db *gorm.DB
}

func (r *urlRepository) GetListUrls(userId uint) ([]*db.URL, error) {
	var urls []*db.URL
	err := r.db.Where("user_id = ?", userId).Find(&urls).Error
	if err != nil {
		return nil, err
	}
	if len(urls) == 0 {
		return []*db.URL{}, nil
	}

	return urls, nil
}

func (r *urlRepository) CreateIfNotExists(url *db.URL) (*db.URL, error) {
	if url == nil {
		return nil, errors.New("URL cannot be nil")
	}

	// Пытаемся вставить новую запись
	result := r.db.Create(url)
	if result.Error == nil {
		return url, nil // Успешная вставка
	}
	var pgErr *pgconn.PgError
	if errors.As(result.Error, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		// Если URL уже существует, находим существующую запись
		var existingURL db.URL
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

func (r *urlRepository) GetBySlug(slug string) (db.URL, error) {
	var url db.URL
	err := r.db.Where("slug_url = ?", slug).First(&url).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return db.URL{}, errors.New("URL not found")
		}
		return db.URL{}, err
	}
	return url, nil
}

func (r *urlRepository) Create(url *db.URL) error {
	if url == nil {
		return errors.New("URL cannot be nil")
	}
	return r.db.Create(url).Error
}

func (r *urlRepository) BatchCreate(urls []db.URL) error {
	if len(urls) == 0 {
		return nil
	}
	return r.db.Create(&urls).Error
}
func AutoMigrateModels(dbUrl string) error {
	models := db.GetModels()
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
