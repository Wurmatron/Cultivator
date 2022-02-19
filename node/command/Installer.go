package command

import (
	"cultivator.wurmatron.io/backend/model"
	"log"
	"strings"
)

func IsBlockchainInstalled(install model.BlockchainInstallation) bool {
	return Exists(install.InstallDir)
}

func InstallBlockchain(install model.BlockchainInstallation) {
	if !Exists("install.sh") {
		RunCommand("wget", install.ScriptDownloadURL)
		RunCommand("ls", "-l")
		RunCommand("unzip", strings.ToLower(install.Name)+".zip")
		RunCommand("ls", "-l")
		RunCommand("rm", strings.ToLower(install.Name)+".zip")
		RunCommand("ls", "-l")
	}
	RunCommand("bash", "./install.sh")
}

func ConfigureNode(chain model.BlockchainInstallation) {
	log.Println("Configuring First time system!")
	RunNodeCommand(chain, "init")
	RunNodeCommand(chain, "keys generate") // TODO Replace with transfer / mnemonic import
	RunNodeCommand(chain, "configure -log-level INFO")
	// TODO Set Log Level To Info
}

func SetupAndRun(install model.BlockchainInstallation, runType string) {
	if IsBlockchainInstalled(install) {
		log.Println(install.Name + " has already been installed!")
	} else {
		log.Println("Starting First Time Setup for '" + install.Name + "' into '" + install.InstallDir + "'")
		InstallBlockchain(install)
		ConfigureNode(install)
	}
	if strings.EqualFold("node", runType) {
		StartFarmer(install)
	}
	if strings.EqualFold("harvester", runType) {
		StartHarvester(install)
	}
}
