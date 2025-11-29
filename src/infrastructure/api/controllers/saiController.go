package controllers

import (
	"fmt"
	"net/http"

	"homelabs-service/src/domain/dtos"
	"homelabs-service/src/domain/entities"
	"homelabs-service/src/domain/queries"
	"homelabs-service/src/infrastructure/repositories"
	"homelabs-service/src/infrastructure/services"
	"homelabs-service/src/shared"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ISAIController struct {
}

func SAIController() *ISAIController {
	return &ISAIController{}
}

func (c *ISAIController) GetItems(ctx fiber.Ctx) error {
	result := shared.ResultData[[]dtos.SAI]()

	items, err := repositories.SAI.GetMany(bson.M{}, 5)
	if err != nil {
		result.AddError(err.Error())
		return ctx.Status(http.StatusInternalServerError).JSON(result.Response())
	}

	var itemList []dtos.SAI
	for _, item := range items {
		itemList = append(itemList, dtos.NewSAI(item))
	}

	result.AddData(itemList)

	return ctx.Status(http.StatusOK).JSON(result.Response())
}

func (c *ISAIController) Create(ctx fiber.Ctx) error {
	result := shared.ResultData[dtos.SAI]()

	bodyData := new(queries.SAI)
	if err := ctx.Bind().Body(bodyData); err != nil {
		result.AddError(fmt.Sprintf("Invalid request body: %s", err.Error()))

		return ctx.Status(http.StatusBadRequest).JSON(result.Response())
	}

	createdData, errors := entities.CreateSAI(*bodyData)
	if len(errors) > 0 {
		for _, err := range errors {
			result.AddError(err)
		}

		return ctx.Status(http.StatusBadRequest).JSON(result.Response())
	}

	item, err := repositories.SAI.Create(*createdData)
	if err != nil {
		result.AddError(err.Error())
		return ctx.Status(http.StatusInternalServerError).JSON(result.Response())
	}

	itemDto := dtos.NewSAI(*item)
	services.SendTelegramSAIMessage(itemDto)

	result.AddData(itemDto)

	return ctx.Status(http.StatusCreated).JSON(result.Response())
}
