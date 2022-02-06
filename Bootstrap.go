package main

import (
	"cultivator.wurmatron.io/node"
	"flag"
	"log"
)

func main() {
	log.SetPrefix("[Bootstrap] > ")
	configurationServer := flag.String("configurationServer", "localhost", "<server ip or domain>")
	log.Println("Connecting to configuration server '" + *configurationServer + "'\n")
	go ConfigureSystem(*configurationServer)
	select {}
}

func ConfigureSystem(configServer string) {
	// TODO Implement
	//backend.Start()
	node.Start()
}
