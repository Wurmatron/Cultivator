package harvester

import (
	"bytes"
	"cultivator.wurmatron.io/backend/model"
	"cultivator.wurmatron.io/node/command"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TODO Temp
var CurrentInstallation = model.BlockchainInstallation{
	Name:              "stai",
	Github:            "https://github.com/STATION-I/stai-blockchain",
	InstallDir:        "stai-blockchain",
	LogFile:           "~/.stai/mainnet/log/debug.log",
	ScriptDownloadURL: "https://raw.githubusercontent.com/Wurmatron/Cultivator/main/scripts/stai.zip",
	Priority:          0,
	ProtectedFiles:    []string{"~/.stai/mainnet/db/blockchain_v1_mainnet.sqlite", "~/.stai/mainnet/wallet/db/"},
}

func Start(configServer string) {
	log.SetPrefix("[Harvester] > ")
	log.Println("Starting up as 'Harvester'")
	// TODO Load Configuration
	log.Println("Running as '" + strings.ToUpper(CurrentInstallation.Name) + "' harvester!")
	command.SetupAndRun(CurrentInstallation, "harvester")
	go ScheduleStatusUpdate(configServer)
}

func ScheduleStatusUpdate(configServer string) {
	for range time.Tick(time.Second * time.Duration(60)) {
		SendHarvesterStatusUpdate(configServer)
	}
}

func SendHarvesterStatusUpdate(server string) {
	plotCount, _ := strconv.ParseInt(command.TotalPlots(CurrentInstallation), 10, 64)
	harvesterStatus := model.Harvester{
		ID:               "",
		Blockchain:       CurrentInstallation.InstallDir,
		LastSync:         time.Now().Unix(),
		Version:          command.Version(CurrentInstallation),
		ConnectionStatus: command.NodeStatus(CurrentInstallation),
		PlotCount:        plotCount,
		PlotSize:         command.TotalPlotSize(CurrentInstallation),
		Drives:           nil,
		PlotHarvestTime:  0,
	}
	log.Println("Sending Harvester Status Update")
	jsonData, e := json.Marshal(harvesterStatus)
	if e != nil {
		log.Println(e.Error())
	}
	req, err := http.NewRequest("POST", server+"/api/metric/harvester", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()
}
