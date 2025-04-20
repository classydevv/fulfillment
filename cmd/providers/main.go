package main

import (
	"log"

	config "github.com/classydevv/fulfillment/configs/providers"
	"github.com/classydevv/fulfillment/internal/providers/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
