package main

import (
	"log"

	"os"

	"github.com/0xERR0R/mailcatcher/config"
	"github.com/0xERR0R/mailcatcher/server"
	"github.com/tkanos/gonfig"
)

func main() {
	var configPath string

	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	configuration := config.Configuration{}
	err := gonfig.GetConf(configPath, &configuration)

	if err != nil {
		log.Fatal("can't read configuration: ", err)
	}

	if err := configuration.Validate(); err != nil {
		log.Fatal("please check the configuration")
	}

	log.Println("Using configuration:", configuration)

	err = server.NewServer(configuration)

	if err != nil {
		log.Fatal("can't start server: ", err)
	}
}
