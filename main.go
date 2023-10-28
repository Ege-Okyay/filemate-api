package main

import (
	"log"

	"github.com/Ege-Okyay/filemate-api/routes"
	"github.com/Ege-Okyay/filemate-api/utils"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	utils.InitDB()

	r := routes.SetupRouter()

	r.Run(":8080")
}
