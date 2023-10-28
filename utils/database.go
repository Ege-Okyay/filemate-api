package utils

import (
	"fmt"
	"os"

	"github.com/Ege-Okyay/filemate-api/models"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	var err error
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&models.User{})
}

func GetDB() *gorm.DB {
	return db
}
