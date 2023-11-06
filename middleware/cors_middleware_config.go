package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://192.168.1.36:8081"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}

	return cors.New(config)
}
