package node

import (
	"cultivator.wurmatron.io/backend/model"
	"cultivator.wurmatron.io/node/command"
	"log"
	"strings"
)

const (
	Blockchain = "chia"
)

// TODO Temp
var CurrentInstallation = model.BlockchainInstallation{
	Name:              "chia",
	Github:            "https://github.com/Chia-Network/chia-blockchain",
	InstallDir:        "chia-blockchain",
	LogFile:           "~/.chia/mainnet/log/debug.log",
	ScriptDownloadURL: "https://raw.githubusercontent.com/Wurmatron/Cultivator/main/scripts/chia.zip",
	Priority:          1,
}

func Start() {
	log.SetPrefix("[Node]      > ")
	log.Println("Starting up as 'Node'")
	// TODO Load Configuration
	log.Println("Running as '" + strings.ToUpper(CurrentInstallation.Name) + "' node!")
	command.SetupAndRun(CurrentInstallation)
}
