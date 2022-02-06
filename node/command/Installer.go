package command

import (
	"cultivator.wurmatron.io/backend/model"
	"log"
	"os"
	"os/exec"
)

func IsBlockchainInstalled(install model.BlockchainInstallation) bool {
	dir, err := os.Stat(install.InstallDir)
	if os.IsNotExist(err) {
		return false
	}
	return dir.IsDir()
}

func InstallBlockchain(install model.BlockchainInstallation) {
	cmd := exec.Command("curl", install.ScriptDownloadURL)
	if output, err := cmd.Output(); err != nil {
		log.Println(output)
	} else {
		log.Fatal(err)
	}
}

func SetupAndRun(install model.BlockchainInstallation) {
	if IsBlockchainInstalled(install) {
		// TODO Run
		log.Println(install.Name + " has already been installed!")
	} else {
		log.Println("Starting First Time Setup for '" + install.Name + "' into '" + install.InstallDir + "'")
		InstallBlockchain(install)
	}
}
