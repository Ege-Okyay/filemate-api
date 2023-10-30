package controllers

import (
	"net/http"

	"github.com/Ege-Okyay/filemate-api/models"
	"github.com/Ege-Okyay/filemate-api/services"
	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing credentials"})
		return
	}

	err := services.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
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
