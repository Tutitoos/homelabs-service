package middlewares

import (
	"net/http"

	"homelabs-service/src/shared"

	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(ctx fiber.Ctx) error {
	result := shared.ResultData[string]()

	url := ctx.OriginalURL()
	if url == "/" || url == "/health" || url == "/monitor" || url == "/openapi.yaml" {
		return ctx.Next()
	}

	authToken := ctx.Get("Authorization", "")
	if authToken == "" {
		result.AddError("Authorization token is missing")

		return ctx.Status(http.StatusUnauthorized).JSON(result.Response())
	}

	if len(authToken) < 7 || authToken[:7] != "Bearer " {
		result.AddError("Authorization token is invalid format")

		return ctx.Status(http.StatusUnauthorized).JSON(result.Response())
	}

	authToken = authToken[7:]
	if authToken != shared.Config.BaseToken {
		result.AddError("Authorization token is invalid")

		return ctx.Status(http.StatusUnauthorized).JSON(result.Response())
	}

	return ctx.Next()
}
