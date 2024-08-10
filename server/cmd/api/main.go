package main

import (
	"log"

	"github.com/pillowskiy/gopix/internal/api"
	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/storage"
)

func main() {
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

	sh := storage.NewStorageHolder(config)
	if err := sh.Setup(); err != nil {
		logger.Fatalf("StorageHolderSetup: %v", err)
	}
	defer sh.Close()

	s := api.NewEchoServer(&config.Server, sh, logger)
	if err := s.Listen(); err != nil {
		logger.Fatalf("ServerListen: %v", err)
	}
}
