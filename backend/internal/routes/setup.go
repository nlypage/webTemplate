package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"webTemplate/internal/api"
)

func Setup(app *fiber.App) {
	app.Use(cors.New())
	apiGroup := app.Group("/api", cors.New())

	apiGroup.Get("/ping", api.Ping)
}
