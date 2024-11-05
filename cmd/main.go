package main

import (
	"webTemplate/cmd/app"
	"webTemplate/internal/adapters/config"
	"webTemplate/internal/adapters/controller/api/setup"
)

func main() {
	appConfig := config.Configure()
	mainApp := app.New(appConfig)

	setup.Setup(mainApp)
	mainApp.Start()
}
