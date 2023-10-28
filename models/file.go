package models

import (
	"github.com/google/uuid"
)

type File struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	FileName string    `json:"fileName"`
	FileSize string    `json:"fileSize"`
	FielPath string    `json:"filePath"`
	UserID   uuid.UUID `json:"userId"`
	User     User      `json:"user" gorm:"foreignKey:UserID"`
}
