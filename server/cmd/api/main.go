package main

import (
	"log"

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

	db, err := storage.NewPostgres(&config.Postgres)
	if err != nil {
		logger.Fatalf("NewPostgres: %v", err)
	} else {
		logger.Infof(
			"Postgres connected with driver '%s', max connections: %v",
			config.Postgres.Driver, db.Stats().MaxOpenConnections,
		)
	}
}
