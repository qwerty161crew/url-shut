package db

import (
	"url-shortener/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MigrateModels(dbUrl string) error {
	models := models.GetMigrationModels()
	db, _ := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}

	return nil
}
