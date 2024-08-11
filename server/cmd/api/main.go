package main

import (
	"log"

	"github.com/pillowskiy/gopix/internal/api"
	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/storage"
)

func main() {
	cfg, err := config.FetchAndLoadConfig()
	if err != nil {
		log.Fatalf("FetchAndLoadConfig: %v", err)
	}

	logger := logger.NewZap(&cfg.Logger).Init()
	logger.Info("Logger successfully initialized")

	sh := storage.NewStorageHolder(cfg)
	if err := sh.Setup(); err != nil {
		logger.Fatalf("StorageHolderSetup: %v", err)
	}
	defer sh.Close()

	s := api.NewEchoServer(&cfg.Server, sh, logger)
	if err := s.Listen(); err != nil {
		logger.Fatalf("ServerListen: %v", err)
	}
}
