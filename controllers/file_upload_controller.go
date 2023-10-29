package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Ege-Okyay/filemate-api/models"
	"github.com/Ege-Okyay/filemate-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context) {
	userID := c.GetString("userId")

	err := c.Request.ParseMultipartForm(10 << 30) // 10 GB
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process form"})
		return
	}

	files := c.Request.MultipartForm.File["files"]

	userFolder := filepath.Join("uploads", userID)
	err = os.MkdirAll(userFolder, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user folder"})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to parse user id"})
		return
	}

	var fileRecords []models.File
	for _, file := range files {
		if len(fileRecords) >= 5 {
			break
		}

		fileID := uuid.New()

		fileName := fmt.Sprintf("%s%s", fileID, filepath.Ext(file.Filename))
		out, err := os.Create(filepath.Join(userFolder, fileName))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a file"})
			return
		}
		defer out.Close()

		filePath := filepath.Join(userFolder, fileName)

		fileObj, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer fileObj.Close()

		err = utils.SaveFileWithBuffer(fileObj, out)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		fileRecord := models.File{
			ID:       fileID,
			FileName: fileName,
			FileSize: file.Size,
			FilePath: filePath,
			UserID:   userUUID,
		}

		fileRecords = append(fileRecords, fileRecord)
	}

	db := utils.GetDB()

	if len(fileRecords) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files to create"})
		return
	}

	for _, record := range fileRecords {
		result := db.Create(&record)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file details"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Files uploaded successfully"})
}

func GetFiles(c *gin.Context) {
	userID := c.GetString("userId")

	db := utils.GetDB()

	var foundFiles []models.File
	result := db.Where("user_id = ?", userID).Find(&foundFiles)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while getting files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": result.Value})
}
