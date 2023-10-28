package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var signingKey = []byte("secret-key")

func GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", err
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return "", err
	}

	return userId, nil
}
