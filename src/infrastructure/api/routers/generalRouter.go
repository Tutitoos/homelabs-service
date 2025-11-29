package routers

import (
	"homelabs-service/src/infrastructure/api/controllers"

	"github.com/gofiber/fiber/v3"
)

func GeneralRouter(router fiber.Router) fiber.Router {
	controllerGeneral := controllers.GeneralController()

	router.Get("/", controllerGeneral.GetHomeHandler)
	router.Get("/health", controllerGeneral.GetHealthHandler)

	return router
}
