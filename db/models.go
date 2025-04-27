package db

import (
	"errors"

	"gorm.io/gorm"
)

func GetModels() []interface{} {
	return []interface{}{User{}, URL{}}
}

var ErrUsernameExists = errors.New("username already exists")

type User struct {
	gorm.Model
	ID       uint   `gorm:"primarykey"`
	Username string `gorm:"type:CHAR(16);unique;not null"`
	Password string `gorm:"not null"`
}
type URL struct {
	gorm.Model
	Slug string `gorm:"type:CHAR(8);unique;not null"`
	Url     string `gorm:"type:TEXT;not null"`
	UserID  uint
	User    User `gorm:"foreignKey:ID"`
}
