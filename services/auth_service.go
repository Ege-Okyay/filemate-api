package services

import (
	"github.com/Ege-Okyay/filemate-api/models"
	"github.com/Ege-Okyay/filemate-api/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.ID = uuid.New()

	db := utils.GetDB()
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

type ReturnError struct {
	Error   error
	Context string
}

func AuthenticateUser(identifier string, password string) (*models.User, ReturnError) {
	db := utils.GetDB()

	var foundUser models.User
	result := db.Where("username = ? OR email = ?", identifier, identifier).First(&foundUser)
	if result.Error != nil {
		return nil, ReturnError{Error: result.Error, Context: "Wrong username or email"}
	}

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if err != nil {
		return nil, ReturnError{Error: err, Context: "Incorrect password"}
	}

	return &foundUser, ReturnError{}
}

func GenerateToken(userID string) (string, error) {
	return utils.GenerateToken(userID)
}
