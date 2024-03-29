package main

import (
	"cultivator.wurmatron.io/backend"
	"flag"
	"log"
	"os"
)

var ConfigurationServer string

func init() {
	flag.StringVar(&ConfigurationServer, "ip", "localhost", "<server ip or domain>")
	flag.Parse()
	if len(os.Getenv("backend_ip")) > 0 {
		ConfigurationServer = os.Getenv("backend_ip")
	}
}

func main() {
	log.SetPrefix("[Bootstrap] > ")
	log.Println("Connecting to configuration server '" + ConfigurationServer + "'\n")
	go ConfigureSystem(ConfigurationServer)
	select {}
}

func ConfigureSystem(configServer string) {
	// TODO Implement
	backend.Start()
	//node.Start(configServer)
	//harvester.Start(configServer)
	//plotting.Start(configServer)
}
