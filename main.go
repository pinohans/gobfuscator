package main

import (
	"flag"
	"log"
)

func main() {
	configFilename := flag.String("c", "", "config.json")
	if *configFilename != "" {
		if err := config.LoadConfig(*configFilename); err != nil {
			log.Fatalln("Failed to LoadConfig: ", err)
		}
	}
	if err := Setup(); err != nil {
		log.Fatalln("Failed to Setup: ", err)
	}

	if err := Obfuscate(); err != nil {
		log.Fatalln("Failed to obfuscate: ", err)
	}
	if err := Build(); err != nil {
		log.Fatalln("Failed to BuildSrc: ", err)
	}
}
