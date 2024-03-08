package main

import (
	"fmt"
	"os"

	"github.com/Ege-Okyay/filemate-api/config"
	"github.com/Ege-Okyay/filemate-api/controllers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/sign-up", controllers.SignUp)
	auth.Post("/login", controllers.Login)
}

func main() {
	config.LoadEnv()

	app := fiber.New()

	config.ConnectDB()
	config.InitFirebase()

	setupRoutes(app)

	app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
