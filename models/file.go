package models

import (
	"github.com/google/uuid"
)

type File struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	FileName string    `json:"fileName"`
	FileSize int64     `json:"fileSize"`
	FilePath string    `json:"filePath"`
	UserID   uuid.UUID `json:"userId"`
	Public   bool      `json:"public"`
}
