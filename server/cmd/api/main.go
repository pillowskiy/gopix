package main

import (
	"fmt"
	"log"

	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/pkg/logger"
)

func main() {
	fmt.Println("Hello World!")

	configFile, err := config.FetchAndLoadConfig()
	if err != nil {
		log.Fatalf("FetchAndLoadConfig: %v", err)
	}

	config, err := config.ParseConfig(configFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	logger := logger.NewZap(&config.Logger).Init()
	logger.Info("Logger successfully initialized")
}
