package main

import (
	"log"
	"os"

	config "github.com/classydevv/fulfillment/configs/providers"
	"github.com/classydevv/fulfillment/internal/providers/app"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env for local deploy. Assumes that the file is called .env
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
