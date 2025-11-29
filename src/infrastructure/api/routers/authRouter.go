package routers

import (
	"homelabs-service/src/infrastructure/api/controllers"

	"github.com/gofiber/fiber/v3"
)

func AuthRouter(router fiber.Router) fiber.Router {
	controllerAuth := controllers.AuthController()

	router.Post("/login", controllerAuth.LoginHandler)

	return router
}
