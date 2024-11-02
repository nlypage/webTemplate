package main

import (
	"webTemplate/cmd/app"
	"webTemplate/internal/adapters/config"
)

func main() {
	appConfig := config.Configure()
	mainApp := app.New(appConfig)

	setup.Setup(mainApp)
	mainApp.Start()
}
