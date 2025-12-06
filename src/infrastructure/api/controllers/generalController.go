package controllers

import (
	"net/http"

	"homelabs-service/src/shared"

	"github.com/gofiber/fiber/v3"
)

type IGeneralController struct {
}

func GeneralController() *IGeneralController {
	return &IGeneralController{}
}

func (c *IGeneralController) GetHomeHandler(ctx fiber.Ctx) error {
	result := shared.ResultData[[]string]()

	result.AddMessage("API is up and running!")

	return ctx.Status(http.StatusOK).JSON(result.Response())
}

func (c *IGeneralController) GetHealthHandler(ctx fiber.Ctx) error {
	result := shared.ResultData[string]()

	result.AddMessage("OK")

	return ctx.Status(http.StatusOK).JSON(result.Response())
}
