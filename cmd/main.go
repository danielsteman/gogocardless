package main

import (
	"fmt"
	"log"

	"github.com/danielsteman/gogocardless/config"
)

func main() {
	config, err := config.LoadAppConfig()
	if err != nil {
		log.Fatal("Couldn't get banks:", err)
	}
	fmt.Println("test:", string(config.SecretID))
}
