package routes

import (
	"github.com/tanayarun/lazydev/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/commit", handlers.GetCommitHandler)
}
