package main

import (
	"cultivator.wurmatron.io/backend"
	"flag"
	"log"
)

func main() {
	log.SetPrefix("[Bootstrap] > ")
	configurationServer := flag.String("configurationServer", "localhost", "<server ip or domain>")
	log.Println("Connecting to configuration server '" + *configurationServer + "'\n")
	go ConfigureSystem()
	select {}
}

func ConfigureSystem() {
	// TODO Implement
	backend.Start()
}
