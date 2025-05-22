package main

import (
	"log"

	"github.com/joho/godotenv"
	config "github.com/classydevv/fulfillment/configs/providers"
	"github.com/classydevv/fulfillment/internal/providers/app"
)

func main() {
	// Load .env for local deploy. Assumes that the file is called .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
