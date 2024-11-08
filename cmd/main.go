package main

import (
	"webTemplate/cmd/app"
	"webTemplate/internal/adapters/config"
	"webTemplate/internal/adapters/controller/api/setup"
)

// @title           WebTemplate API
// @version         1.0
// @description     This is a webtemplate API that contains project dir structure, JWT auth, basic user entitites and can be further expanded.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description "Type 'Bearer TOKEN' to correctly set the API Key"
func main() {
	appConfig := config.Configure()
	mainApp := app.New(appConfig)

	setup.Setup(mainApp)
	mainApp.Start()
}
