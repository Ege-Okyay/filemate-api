package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Ege-Okyay/filemate-api/models"
	"github.com/Ege-Okyay/filemate-api/services"
	"github.com/Ege-Okyay/filemate-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context) {
	userID := c.GetString("userId")
	var err error

	err = c.Request.ParseMultipartForm(10 << 30) // 10 GB
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
			Public:   false,
		}

		fileRecords = append(fileRecords, fileRecord)
	}

	if len(fileRecords) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files to create"})
		return
	}

	for _, record := range fileRecords {
		err = services.CreateFileRecord(record)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file details"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Files uploaded successfully"})
}

func DeleteFile(c *gin.Context) {
	var err error

	fileID := c.Query("fileId")

	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file id"})
		return
	}

	foundFile, err := services.FindFileByID(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	err = os.Remove(foundFile.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting file from server"})
		return
	}

	err = services.DeleteFileRecord(foundFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted the file"})
}

func DownloadFile(c *gin.Context) {
	fileID := c.Query("fileId")

	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file id"})
		return
	}

	foundFile, err := services.FindFileByID(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	filePath := foundFile.FilePath

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Dispoistion", "attachment filename="+foundFile.FileName)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer_Encoding", "binary")
	c.Header("Expires", "0")

	c.File(filePath)
}

func ChangeFilePublicty(c *gin.Context) {
	var err error

	fileID := c.Query("fileId")
	newStatus := c.Query("public")

	if fileID == "" || newStatus == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query data"})
		return
	}

	newStatusValue, err := strconv.ParseBool(newStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Status value parsing error"})
		return
	}

	foundFile, err := services.FindFileByID(fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File not found"})
		return
	}

	err = services.UpdateFileStatus(foundFile, newStatusValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating file status"})
		return
	}

	if newStatusValue == false {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully made the file private"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully made the file public"})
		return
	}
}

func GetFiles(c *gin.Context) {
	userID := c.GetString("userId")

	foundFiles, err := services.FindFilesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while reading files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": foundFiles})
}
