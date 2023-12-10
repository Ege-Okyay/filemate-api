package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Ege-Okyay/filemate-api/controllers"
	"github.com/Ege-Okyay/filemate-api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CorsMiddleware())

	authMiddleware := middleware.AuthMiddleware()

	r.GET("/health-check", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{})
		return
	})

	r.POST("/post-check", func(ctx *gin.Context) {
		var reqBody map[string]interface{}

		err := json.NewDecoder(ctx.Request.Body).Decode(&reqBody)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"body": reqBody})
		return
	})

	auth := r.Group("/auth")
	{
		auth.POST("/sign-up", controllers.SignUp)
		auth.POST("/login", controllers.Login)
	}

	file := r.Group("/file")
	file.Use(authMiddleware.MiddlewareFunc())
	{
		file.POST("/upload", controllers.UploadFile)
		file.POST("/change-publicty", controllers.ChangeFilePublicty)

		file.GET("/files", controllers.GetFiles)
		file.GET("/download", controllers.DownloadFile)
		file.DELETE("/delete", controllers.DeleteFile)
	}

	return r
}
