package routes

import (
	"github.com/Ege-Okyay/filemate-api/controllers"
	"github.com/Ege-Okyay/filemate-api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	authMiddleware := middleware.AuthMiddleware()

	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", controllers.SignUp)
		auth.POST("/login", controllers.Login)
	}

	private := r.Group("/private")
	private.Use(authMiddleware.MiddlewareFunc())
	{
		private.GET("/user", controllers.UserProfile)
	}

	return r
}
