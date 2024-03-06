package config

import (
	"fmt"
	"os"

	"github.com/Ege-Okyay/filemate-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"postgres://default:%s@ep-green-mode-a2eum8tp.eu-central-1.aws.neon.tech:5432/verceldb?sslmode=require",
		os.Getenv("DB_SECRET"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	DB = db

	DB.AutoMigrate(&models.User{})
}
