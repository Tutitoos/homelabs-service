package routers

import (
	"homelabs-service/src/infrastructure/api/controllers"

	"github.com/gofiber/fiber/v3"
)

func TelegramRouter(router fiber.Router) fiber.Router {
	controller := controllers.TelegramController()

	router.Post("/sai", controller.SAI)
	router.Post("/dns", controller.DNS)
	router.Post("/backups", controller.Backup)

	return router
}
