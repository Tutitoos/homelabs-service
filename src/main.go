package main

import (
	"log"

	"homelabs-service/src/infrastructure/api"
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

	shared.Logger.Infof("Starting Homelabs Service in %s mode", shared.Config.AppEnv)

	api := api.Api()
	api.Start()
}
