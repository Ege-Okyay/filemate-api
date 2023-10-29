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

	file := r.Group("/file")
	file.Use(authMiddleware.MiddlewareFunc())
	{
		file.POST("/upload", controllers.UploadFile)
		file.GET("/files", controllers.GetFiles)
		file.GET("/download", controllers.DownloadFile)
		file.DELETE("/delete", controllers.DeleteFile)
	}

	return r
}
