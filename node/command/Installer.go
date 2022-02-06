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
	if !Exists("scripts/install.sh") {
		RunCommand("wget", install.ScriptDownloadURL)
		RunCommand("unzip", strings.ToLower(install.Name)+".zip")
		RunCommand("rm", strings.ToLower(install.Name)+".zip")
	}
	RunCommand("bash", "scripts/install.sh")
}

func ConfigureNode(chain model.BlockchainInstallation) {
	log.Println("Configuring First time system!")
	RunNodeCommand(chain, "init")
	RunNodeCommand(chain, "keys generate") // TODO Replace with transfer / mnemonic import
	RunNodeCommand(chain, "configure -log-level INFO")
}

func SetupAndRun(install model.BlockchainInstallation) {
	if IsBlockchainInstalled(install) {
		log.Println(install.Name + " has already been installed!")
	} else {
		log.Println("Starting First Time Setup for '" + install.Name + "' into '" + install.InstallDir + "'")
		InstallBlockchain(install)
		ConfigureNode(install)
	}
	StartNode(install)
}
