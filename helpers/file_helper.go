package helpers

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Ege-Okyay/filemate-api/config"
	"github.com/Ege-Okyay/filemate-api/models"
	"github.com/google/uuid"
)

func UploadAndSaveFile(ctx context.Context, file io.Reader, filename string, userID string) error {
	client, err := config.FirebaseApp.Storage(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Firebase Storage client: %v", err)
	}

	bucketName := fmt.Sprintf("%s.appspot.com", os.Getenv("FIREBASE_PROJECT_ID"))
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return fmt.Errorf("failed to initialize Firebase Bucket: %v", err)
	}

	objectPath := fmt.Sprintf("%s/%s", userID, uuid.New().String())
	wc := bucket.Object(objectPath).NewWriter(ctx)

	buffer := make([]byte, 1024*1024)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			wc.Close()
			return fmt.Errorf("failed to read from file: %v", err)
		}
		if n == 0 {
			break
		}

		if _, err := wc.Write(buffer[:n]); err != nil {
			wc.Close()
			return fmt.Errorf("failed to write to Firebase Storage: %v", err)
		}
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close Firebase Storage writer: %v", err)
	}

	attrs, err := bucket.Object(objectPath).Attrs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get object attributes: %v", err)
	}

	url := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media", bucketName, attrs.Name)

	fileModel := models.File{
		Name:      filename,
		Path:      url,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	dbCtx := context.TODO()
	_, err = config.FileCollection.InsertOne(dbCtx, fileModel)
	if err != nil {
		return fmt.Errorf("failed to save file to database: %v", err)
	}

	return nil
}
