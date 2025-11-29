package main

import (
	"log"

	"homelabs-service/src/infrastructure/api"
	"homelabs-service/src/infrastructure/datasources"
	"homelabs-service/src/infrastructure/repositories"
	"homelabs-service/src/shared"
)

func main() {
	shared.CapturePanic()

	err := shared.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		return
	}

	shared.InitLogger()

	dataSources := datasources.DataSources()
	dataSources.Connect()

	repositories.InitializeRepositories()

	shared.Logger.Infof("Starting Database Service in %s mode", shared.Config.AppEnv)

	api := api.Api()
	api.Start()
}
