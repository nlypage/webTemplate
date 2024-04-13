package main

import (
	"github.com/gofiber/fiber/v3"
	"os"
	"webTemplate/internal/database"
	"webTemplate/internal/routes"
)

func main() {
	database.Connect()
	port := os.Getenv("BACKEND_PORT")
	app := fiber.New()
	routes.Setup(app)

	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}
