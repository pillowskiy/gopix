package main

import (
	"fmt"
	"log"

	"github.com/pillowskiy/gopix/internal/config"
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

	fmt.Println(config)
}
