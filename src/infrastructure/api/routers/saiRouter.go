package routers

import (
	"homelabs-service/src/infrastructure/api/controllers"

	"github.com/gofiber/fiber/v3"
)

func SAIRouter(router fiber.Router) fiber.Router {
	controller := controllers.SAIController()

	r := router.Group("/sai")
	r.Get("/", controller.GetItems)
	r.Post("/", controller.Create)

	return router
}
