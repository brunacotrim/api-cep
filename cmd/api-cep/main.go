package main

import (
	"log"

	"api-cep/config"
	"api-cep/service"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	service, err := service.New(config)
	if err != nil {
		log.Fatalf("error create service: %v", err)
	}

	service.StartServer()

}
