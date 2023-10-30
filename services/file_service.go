package services

import (
	"github.com/Ege-Okyay/filemate-api/models"
	"github.com/Ege-Okyay/filemate-api/utils"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	db = utils.GetDB()
}

func CreateFileRecords(fileRecords []models.File) error {
	for _, record := range fileRecords {
		result := db.Create(&record)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func FindFileByID(fileID string) (*models.File, error) {
	var foundFile models.File
	result := db.Where("id = ?", fileID).First(&foundFile)
	if result.Error != nil {
		return nil, result.Error
	}

	return &foundFile, nil
}

func DeleteFileRecord(file *models.File) error {
	deleteResult := db.Delete(file)
	if deleteResult.Error != nil {
		return deleteResult.Error
	}

	return nil
}

func UpdateFileStatus(file *models.File, public bool) error {
	err := db.Model(file).Update("public", public)
	if err.Error != nil {
		return err.Error
	}

	return nil
}

func FindFilesByUserID(userID string) ([]models.File, error) {
	var foundFiles []models.File
	result := db.Where("user_id = ?", userID).Find(&foundFiles)
	if result.Error != nil {
		return nil, result.Error
	}

	return foundFiles, nil
}
