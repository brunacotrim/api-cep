package main

import (
	"fmt"
	"log"

	"api-cep/config"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatalf("error loagind configuration: %v", err)
	}

	fmt.Println(config)
}
