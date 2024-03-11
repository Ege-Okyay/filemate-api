package controllers

import (
	"bytes"
	"io/ioutil"

	"github.com/Ege-Okyay/filemate-api/helpers"
	"github.com/gofiber/fiber/v2"
)

func UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	defer src.Close()

	data, err := ioutil.ReadAll(src)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	url, err := helpers.UploadFile(c.Context(), bytes.NewReader(data), file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"url": url,
	})
}
