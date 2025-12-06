package controllers

import (
	"fmt"
	"net/http"

	"homelabs-service/src/domain/dtos"
	"homelabs-service/src/domain/queries"
	"homelabs-service/src/infrastructure/services"
	"homelabs-service/src/shared"

	"github.com/gofiber/fiber/v3"
)

type ITelegramController struct {
}

func TelegramController() *ITelegramController {
	return &ITelegramController{}
}

func (c *ITelegramController) SAI(ctx fiber.Ctx) error {
	result := shared.ResultData[dtos.SAI]()

	bodyData := new(queries.SAI)
	if err := ctx.Bind().Body(bodyData); err != nil {
		result.AddError(fmt.Sprintf("Invalid request body: %s", err.Error()))

		return ctx.Status(http.StatusBadRequest).JSON(result.Response())
	}

	itemDto := dtos.NewSAI(*bodyData)
	services.SendTelegramSAIMessage(itemDto)

	result.AddData(itemDto)

	return ctx.Status(http.StatusCreated).JSON(result.Response())
}

func (c *ITelegramController) DNS(ctx fiber.Ctx) error {
	result := shared.ResultData[dtos.DNS]()

	bodyData := new(queries.DNS)
	if err := ctx.Bind().Body(bodyData); err != nil || bodyData.DNSId == nil || bodyData.StatusId == nil {
		bodyData = shared.PARSER.ParseFormData(string(ctx.Body()))

		if bodyData.DNSId == nil || bodyData.StatusId == nil {
			shared.Logger.Error("DNSController.Create: Missing required fields")
			result.AddError("dns_id and status_id are required")
			return ctx.Status(http.StatusBadRequest).JSON(result.Response())
		}
	}

	itemDto := dtos.NewDNS(*bodyData)
	services.SendTelegramDNSMessage(itemDto)

	result.AddData(itemDto)

	return ctx.Status(http.StatusCreated).JSON(result.Response())
}

func (c *ITelegramController) Backup(ctx fiber.Ctx) error {
	result := shared.ResultData[dtos.Backup]()

	bodyData := new(queries.Backup)
	if err := ctx.Bind().Body(bodyData); err != nil {
		result.AddError(fmt.Sprintf("Invalid request body: %s", err.Error()))

		return ctx.Status(http.StatusBadRequest).JSON(result.Response())
	}

	itemDto := dtos.NewBackup(*bodyData)
	services.SendTelegramBackupMessage(itemDto)

	result.AddData(itemDto)

	return ctx.Status(http.StatusCreated).JSON(result.Response())
}
