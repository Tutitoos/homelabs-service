package routers

import (
	"homelabs-service/src/infrastructure/api/controllers"

	"github.com/gofiber/fiber/v3"
)

func DNSRouter(router fiber.Router) fiber.Router {
	controller := controllers.DNSController()

	router.Get("/", controller.GetItems)
	router.Post("/", controller.Create)

	return router
}
