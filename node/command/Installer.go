package command

import (
	"cultivator.wurmatron.io/backend/model"
	"log"
	"os"
	"os/exec"
	"strings"
)

func IsBlockchainInstalled(install model.BlockchainInstallation) bool {
	return exists(install.InstallDir)
}

func exists(d string) bool {
	dir, err := os.Stat(d)
	if os.IsNotExist(err) {
		return false
	}
	return dir.IsDir()
}

func InstallBlockchain(install model.BlockchainInstallation) {
	if !exists("scripts/install.sh") {
		RunCommand("wget", install.ScriptDownloadURL)
		RunCommand("unzip", strings.ToLower(install.Name)+".zip")
		RunCommand("rm", strings.ToLower(install.Name)+".zip")
	}
	RunCommand("bash", "scripts/install.sh")
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

func RunCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Println("Failed to run command '" + name + strings.Join(arg, " ") + "'")
		log.Println(err.Error())
	}
}
