package models

import "gorm.io/gorm"

const (
	TypeSimpleUtterance = "SimpleUtterance"
)

// Request описывает запрос пользователя.
// см. https://yandex.ru/dev/dialogs/alice/doc/request.html
type RequestCreateUrl struct {
	Url string `json:"url"`
}

type ResponseCreateUrl struct {
	Result string `json:"result"`
}

type URL struct {
	gorm.Model
	SlugUrl string `gorm:"type:CHAR(8);unique;not null"`
	Url     string `gorm:"type:TEXT;not null"`
}

func GetMigrationModels() []interface{} {
	return []interface{}{
		&URL{},
	}
}
func (u *URL) GetByID(db *gorm.DB, slug string) (URL, error) {
	var url URL
	result := db.Where("slug_url = ?", slug).First(&url)
	if result.Error != nil {
		return URL{}, result.Error
	}
	return url, nil
}

func (u *URL) Create(db *gorm.DB, url *URL) error {
	result := db.Create(url)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
