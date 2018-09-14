package main

import (
	"log"

	"github.com/tkanos/gonfig"
	"os"
)

func main() {

	var configPath string
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	configuration := Configuration{}
	err := gonfig.GetConf(configPath, &configuration)
	if err != nil {
		log.Fatal("can't read configuration: ", err)
	}

	if err := configuration.Validate(); err != nil {

		log.Fatal("please check the configuration")
	}

	log.Println("Using configuration:", configuration)

	err1 := NewServer(&configuration)
	if err1 != nil {
		log.Fatal("can't start server: ", err)
	}
}
