package services

import (
	"github.com/Ege-Okyay/filemate-api/models"
	"github.com/Ege-Okyay/filemate-api/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ReturnError represents an error object to be returned from the service functions.
type ReturnError struct {
	Error   error
	Context string
}

// CheckForCredentials checks if the username or email is already taken in the database.
// It returns an error object and two boolean values indicating if the username or email is taken.
func CheckForTakenCredentials(user *models.User) (ReturnError, bool, bool) {
	db := utils.GetDB()

	var userCount, emailCount int64
	var err error

	// Check the count of users with the given username in the database.
	err = db.Model(&models.User{}).Where("username = ?", user.Username).Count(&userCount).Error
	if err != nil {
		return ReturnError{Error: err, Context: "Failed to search username in the database"}, false, false
	}

	// Check the count of users with the given email in the database.
	err = db.Model(&models.User{}).Where("email = ?", user.Email).Count(&emailCount).Error
	if err != nil {
		return ReturnError{Error: err, Context: "Failed to search email in the database"}, false, false
	}

	return ReturnError{}, userCount > 0, emailCount > 0
}

// CreateUser creates a new user in the database after hashing the password.
// It returns an error object if there is an issue during the creation process.
func CreateUser(user *models.User) ReturnError {
	// Hash the user's password for security.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return ReturnError{Error: err, Context: "Failed to hash the password"}
	}

	// Store the hashed password and generate a UUID for the user.
	user.Password = string(hashedPassword)
	user.ID = uuid.New()

	db := utils.GetDB()
	result := db.Create(&user)
	if result.Error != nil {
		return ReturnError{Error: result.Error, Context: "Failed to create user"}
	}

	return ReturnError{}
}

// AuthenticateUser checks if the given identifier and password match an existing user in the database.
// It returns the user object and an error if the authentication fails.
func AuthenticateUser(identifier string, password string) (*models.User, ReturnError) {
	db := utils.GetDB()

	var foundUser models.User
	result := db.Where("username = ? OR email = ?", identifier, identifier).First(&foundUser)
	if result.Error != nil {
		return nil, ReturnError{Error: result.Error, Context: "Wrong username or email"}
	}

	// Compare the provided password with the hashed password in the database.
	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if err != nil {
		return nil, ReturnError{Error: err, Context: "Wrong password"}
	}

	return &foundUser, ReturnError{}
}

// GenerateToken generates a token for the provided user ID.
// It returns a string representation of the token and an error if token generation fails.
func GenerateToken(userID string) (string, error) {
	return utils.GenerateToken(userID)
}
