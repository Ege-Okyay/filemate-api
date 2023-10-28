package middleware

import (
	"net/http"

	"github.com/Ege-Okyay/filemate-api/utils"
	"github.com/gin-gonic/gin"
)

type authMiddleware struct {
}

func AuthMiddleware() authMiddleware {
	return authMiddleware{}
}

func (a *authMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not found"})
			c.Abort()
			return
		}

		userId, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
