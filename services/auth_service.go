package services

import (
	"github.com/Ege-Okyay/filemate-api/config"
	"github.com/Ege-Okyay/filemate-api/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(username string, email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{Username: username, Email: email, Password: string(hashedPassword)}
	if err := config.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func Login(email string, password string) (string, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(config.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
