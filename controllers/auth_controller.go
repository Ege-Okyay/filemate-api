package controllers

import (
	"net/http"

	"github.com/Ege-Okyay/filemate-api/models"
	"github.com/Ege-Okyay/filemate-api/services"
	"github.com/gin-gonic/gin"
)

// SignUp signs user in if there aren't any errors while validating credentials.
// It takes user model fields from the POST request and returns a JSON with the suitable HTTP code and a message or error.
func SignUp(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing credentials"})
		return
	}

	checkCredentials, usernameFound, emailFound := services.CheckForTakenCredentials(&user)
	if checkCredentials.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": checkCredentials.Context})
		return
	}

	if usernameFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	if emailFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	createUser := services.CreateUser(&user)
	if createUser.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": createUser.Context})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully signed up"})
}

func Login(c *gin.Context) {
	var loginData struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing credentials"})
		return
	}

	foundUser, err := services.AuthenticateUser(loginData.Identifier, loginData.Password)
	if err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Context})
		return
	}

	var tokenErr error
	token, tokenErr := services.GenerateToken(foundUser.ID.String())
	if tokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in", "token": token})
}

func UserProfile(c *gin.Context) {
	userID := c.GetString("userId")
	c.JSON(http.StatusOK, gin.H{"userID": userID})
}
