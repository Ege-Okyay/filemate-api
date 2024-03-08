package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

var UserCollection *mongo.Collection
var FileCollection *mongo.Collection

func ConnectDB() {
	dsn := fmt.Sprintf(
		"mongodb+srv://%s:%s@cluster0.fvfflwp.mongodb.net/%s?retryWrites=true&w=majority&appName=Cluster0",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	clientOptions := options.Client().ApplyURI(dsn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	userCollection := client.Database("Filemate").Collection("users")
	fileCollection := client.Database("Filemate").Collection("files")

	UserCollection = userCollection
	FileCollection = fileCollection

	Client = client
}
