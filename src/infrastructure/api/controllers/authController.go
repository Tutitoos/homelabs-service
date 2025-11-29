package controllers

import (
	"net/http"

	"homelabs-service/src/shared"
	"github.com/gofiber/fiber/v3"
)

type IAuthController struct {
}

func AuthController() *IAuthController {
	return &IAuthController{}
}

func (c *IAuthController) LoginHandler(ctx fiber.Ctx) error {
	result := shared.ResultData[string]()

	result.AddMessage("Login successful!")

	return ctx.Status(http.StatusOK).JSON(result.Response())
}
