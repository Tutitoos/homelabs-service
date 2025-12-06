package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"homelabs-service/src/infrastructure/api/middlewares"
	"homelabs-service/src/infrastructure/api/routers"
	"homelabs-service/src/shared"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type ApiStructure interface {
	CreateApp() *fiber.App
	Start()
}

type ApiService struct {
}

func Api() ApiStructure {
	return &ApiService{}
}

func (a *ApiService) CreateApp() *fiber.App {
	defer shared.CapturePanic()

	shared.Collector()

	app := fiber.New(fiber.Config{
		AppName:      "Homelabs Service",
		ServerHeader: "Fiber",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			result := shared.ResultData[string]()
			result.AddError(err.Error())
			fmt.Println("Internal Server Error: " + err.Error())
			return ctx.Status(http.StatusInternalServerError).JSON(result.Response())
		},
	})

	app.Use(recover.New())
	app.Use(helmet.New())

	app.Use(logger.New(logger.Config{
		Format:      "${pid} :: ${time} :: ${ip} :: ${ips} :: ${status} :: ${method} ${path} :: ${latency}\n",
		TimeFormat:  "02-01-2006 03:04:05 PM",
		TimeZone:    "Europe/Madrid",
		ForceColors: true,
		Next: func(ctx fiber.Ctx) bool {
			url := ctx.OriginalURL()

			return url == "/monitor" || strings.Contains(url, "/swagger") || ctx.Method() == "OPTIONS"
		},
	}))

	allowOrigins := shared.Config.BaseAllowedOrigins
	allowCredentials := false
	if len(allowOrigins) == 1 && allowOrigins[0] != "*" {
		allowCredentials = true
	}

	app.Use(cors.New(cors.Config{
		AllowCredentials: allowCredentials,
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{fiber.MethodOptions, fiber.MethodGet, fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
	}))

	routers.GeneralRouter(app)

	app.Use(middlewares.AuthMiddleware)

	routers.AuthRouter(app.Group("/auth"))
	routers.TelegramRouter(app.Group("/telegram"))

	return app
}

func (a *ApiService) Start() {
	app := a.CreateApp()
	err := app.Listen(fmt.Sprintf(":%d", shared.Config.BasePort), fiber.ListenConfig{
		EnablePrefork: shared.Config.BaseMultiProcess,
		TLSMinVersion: tls.VersionTLS12,
	})
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}
