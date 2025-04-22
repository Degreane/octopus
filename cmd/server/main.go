package main

import (
	"github.com/degreane/octopus/internal/service/logger"
	"github.com/joho/godotenv"
)

func main() {
	// initialize Logger
	log := logger.GetLogger().WithField("component", "Server")
	log.Info("Starting Octopus Server")
	// load  .env file
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file", err)
	}
}
