package main

import (
	"fmt"
	"os"

	"github.com/Ege-Okyay/filemate-api/config"
	"github.com/Ege-Okyay/filemate-api/controllers"
	"github.com/Ege-Okyay/filemate-api/middleware"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/sign-up", controllers.SignUp)
	auth.Post("/login", controllers.Login)

	file := app.Group("/file")
	file.Use(middleware.JWTProtected())
	file.Post("/upload", controllers.UploadFile)
}

func main() {
	config.LoadEnv()

	app := fiber.New()

	config.ConnectDB()
	config.InitFirebase()

	setupRoutes(app)

	app.Listen(fmt.Sprintf("localhost:%s", os.Getenv("PORT")))
}
