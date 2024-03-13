package services

import (
	"context"
	"time"

	"github.com/Ege-Okyay/filemate-api/config"
	"github.com/Ege-Okyay/filemate-api/models"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(username string, email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{Username: username, Email: email, Password: string(hashedPassword), CreatedAt: time.Now(), UpdatedAt: time.Now()}

	ctx := context.TODO()
	_, err = config.UserCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func Login(email string, password string) (string, error) {
	var user models.User

	ctx := context.TODO()
	filter := bson.M{"email": email}

	err := config.UserCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = user
	tokenString, err := token.SignedString(config.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
