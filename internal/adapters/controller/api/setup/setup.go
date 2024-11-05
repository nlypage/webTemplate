package setup

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
	"webTemplate/cmd/app"
)

func Setup(app *app.App) {
	app.Fiber.Use(cors.New(cors.ConfigDefault))

	if viper.GetBool("settings.debug") {
		app.Fiber.Use(logger.New(logger.Config{TimeZone: viper.GetString("settings.timezone")}))
	}

	app.Fiber.Get("/ping", func(c *fiber.Ctx) error {
		c.Status(fiber.StatusOK)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": true,
			"body":   "pong",
		})
	})

	/*
		// Setup api v1 routes
		apiV1 := app.Fiber.Group("/api/v1")

		middlewareHandler := middlewares.NewMiddlewareHandler(app)

		// Setup user routes
		userHandler := v1.NewUserHandler(app)
		userHandler.Setup(apiV1, middlewareHandler.IsAuthenticated)

		// Setup question routes
		questionHandler := v1.NewQuestionHandler(app)
		questionHandler.Setup(apiV1, middlewareHandler.IsAuthenticated)

		// Setup conference routes
		conferenceHandler := v1.NewConferenceHandler(app)
		conferenceHandler.Setup(apiV1, middlewareHandler.IsAuthenticated)

		eventHandler := v1.NewEventHandler(app)
		eventHandler.Setup(apiV1, middlewareHandler.IsAuthenticated)
	*/
}
