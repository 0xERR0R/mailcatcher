package main

import (
	"log"

	"github.com/0xERR0R/mailcatcher/config"
	"github.com/0xERR0R/mailcatcher/server"
)

func main() {
	configuration, err := config.GetConfiguration()

	if err != nil {
		log.Fatal("configuration error: ", err)
	}

	log.Println("Using configuration:", configuration)

	err = server.NewServer(configuration)

	if err != nil {
		log.Fatal("can't start server: ", err)
	}
}
