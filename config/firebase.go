package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebase() {
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS_PATH"))
	config := &firebase.Config{ProjectID: os.Getenv("FIREBASE_PROJECT_ID")}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing firebase application: %v\n", err)
	}

	FirebaseApp = app
}
