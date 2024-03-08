package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebase() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("path")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatal("Error initializing Firebase app: %v\n", err)
	}

	FirebaseApp = app
}
