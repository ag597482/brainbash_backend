package main

import (
	"log"

	"brainbash_backend/config"
	"brainbash_backend/internal/app"
	appMongo "brainbash_backend/internal/mongo"
)

var appConfig config.AppConfig

func main() {
	config.InitGlobalConfig(&appConfig)

	appMongo.Init(&appConfig)

	application, err := app.NewApp(&appConfig)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	if err := application.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
