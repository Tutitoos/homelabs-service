package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"homelabs-service/src/domain/dtos"
	"homelabs-service/src/domain/entities"
	"homelabs-service/src/domain/queries"
	"homelabs-service/src/infrastructure/repositories"
	"homelabs-service/src/infrastructure/services"
	"homelabs-service/src/shared"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IDNSController struct {
}

func DNSController() *IDNSController {
	return &IDNSController{}
}

func (c *IDNSController) GetItems(ctx fiber.Ctx) error {
	result := shared.ResultData[[]dtos.DNS]()

	items, err := repositories.DNS.GetMany(bson.M{}, 5)
	if err != nil {
		result.AddError(err.Error())
		return ctx.Status(http.StatusInternalServerError).JSON(result.Response())
	}

	var itemList []dtos.DNS
	for _, item := range items {
		itemList = append(itemList, dtos.NewDNS(item))
	}

	result.AddData(itemList)

	return ctx.Status(http.StatusOK).JSON(result.Response())
}

func (c *IDNSController) Create(ctx fiber.Ctx) error {
	result := shared.ResultData[dtos.DNS]()

	bodyData := new(queries.DNS)
	if err := ctx.Bind().Body(bodyData); err != nil || bodyData.DNSId == nil || bodyData.StatusId == nil {
		bodyData = c.parseFormData(string(ctx.Body()))

		if bodyData.DNSId == nil || bodyData.StatusId == nil {
			shared.Logger.Error("DNSController.Create: Missing required fields")
			result.AddError("dns_id and status_id are required")
			return ctx.Status(http.StatusBadRequest).JSON(result.Response())
		}
	}

	createdData, errors := entities.CreateDNS(*bodyData)
	if len(errors) > 0 {
		for _, err := range errors {
			result.AddError(err)
		}

		shared.Logger.Error("DNSController.Create: Validation errors", "errors", errors)

		return ctx.Status(http.StatusBadRequest).JSON(result.Response())
	}

	item, err := repositories.DNS.Create(*createdData)
	if err != nil {
		shared.Logger.Error("DNSController.Create: Error creating DNS", "error", err)
		result.AddError(err.Error())
		return ctx.Status(http.StatusInternalServerError).JSON(result.Response())
	}

	itemDto := dtos.NewDNS(*item)
	services.SendTelegramDNSMessage(itemDto)

	result.AddData(itemDto)

	return ctx.Status(http.StatusCreated).JSON(result.Response())
}

func (c *IDNSController) parseFormData(rawBody string) *queries.DNS {
	bodyData := new(queries.DNS)

	for _, pair := range strings.Split(rawBody, "&") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "dns_id":
			if dnsId, err := strconv.Atoi(val); err == nil {
				bodyData.DNSId = &dnsId
			}
		case "status_id":
			if statusId, err := strconv.Atoi(val); err == nil {
				bodyData.StatusId = &statusId
			}
		case "created_at":
			if createdAt, err := strconv.ParseInt(val, 10, 64); err == nil {
				bodyData.CreatedAt = &createdAt
			}
		}
	}

	return bodyData
}
