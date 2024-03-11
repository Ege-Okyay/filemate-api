package helpers

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Ege-Okyay/filemate-api/config"
)

func UploadFile(ctx context.Context, file io.Reader, filename string) (string, error) {
	client, err := config.FirebaseApp.Storage(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create Firebase Storage client: %v", err)
	}

	bucketName := fmt.Sprintf("%s.appspot.com", os.Getenv("FIREBASE_PROJECT_ID"))
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to initialize Firebase Bucket: %v", err)
	}

	objectPath := filename
	wc := bucket.Object(objectPath).NewWriter(ctx)

	buffer := make([]byte, 1024*1024)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			wc.Close()
			return "", fmt.Errorf("failed to read from file: %v", err)
		}
		if n == 0 {
			break
		}

		if _, err := wc.Write(buffer[:n]); err != nil {
			wc.Close()
			return "", fmt.Errorf("failed to write to Firebase Storage: %v", err)
		}
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close Firebase Storage writer: %v", err)
	}

	url := fmt.Sprintf("gs://%s/%s", os.Getenv("FIREBASE_STORAGE_BUCKET"), filename)
	return url, nil
}
