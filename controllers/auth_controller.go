package controllers

import (
	"github.com/Ege-Okyay/filemate-api/services"
	"github.com/gofiber/fiber/v2"
)

func SignUp(c *fiber.Ctx) error {
	var requestData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	if err := services.SignUp(requestData.Username, requestData.Email, requestData.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Successfully signed up"})
}

func Login(c *fiber.Ctx) error {
	var requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	token, err := services.Login(requestData.Email, requestData.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid username or password"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}
