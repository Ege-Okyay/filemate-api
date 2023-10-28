package controllers

import (
	"net/http"

	"github.com/Ege-Okyay/filemate-api/models"
	"github.com/Ege-Okyay/filemate-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing credentials"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": " Failed to process password"})
		return
	}

	user.Password = string(hashedPassword)
	user.ID = uuid.New()

	db := utils.GetDB()
	result := db.Create(&user)
	if result.Error != nil {
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

	db := utils.GetDB()

	var foundUser models.User
	result := db.Where("username = ? OR email = ?", loginData.Identifier, loginData.Identifier).First(&foundUser)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wrong username or email"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	token, err := utils.GenerateToken(foundUser.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in", "token": token})
}

func UserProfile(c *gin.Context) {
	userID := c.GetString("userId")
	c.JSON(http.StatusOK, gin.H{"userID": userID})
}
