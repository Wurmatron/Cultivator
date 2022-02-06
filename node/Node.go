package node

import (
	"bytes"
	"cultivator.wurmatron.io/backend/model"
	"cultivator.wurmatron.io/node/command"
	"encoding/json"
	"fmt"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	Blockchain = "chia"
)

// TODO Temp
var CurrentInstallation = model.BlockchainInstallation{
	Name:              "chia",
	Github:            "https://github.com/Chia-Network/chia-blockchain",
	InstallDir:        "/home/wurmatron/Cultivator/chia-blockchain",
	LogFile:           "~/.chia/mainnet/log/debug.log",
	ScriptDownloadURL: "https://raw.githubusercontent.com/Wurmatron/Cultivator/main/scripts/chia.zip",
	Priority:          1,
}

func Start(configServer string) {
	log.SetPrefix("[Node]      > ")
	log.Println("Starting up as 'Node'")
	// TODO Load Configuration
	log.Println("Running as '" + strings.ToUpper(CurrentInstallation.Name) + "' node!")
	command.SetupAndRun(CurrentInstallation)
	go ScheduleStatusUpdate(configServer)
}

func ScheduleStatusUpdate(configServer string) {
	for range time.Tick(time.Second * time.Duration(60)) {
		SendNodeStatusUpdate(configServer)
	}
}

func SendNodeStatusUpdate(server string) {
	plotCount, _ := strconv.ParseInt(command.TotalPlots(CurrentInstallation), 10, 64)
	mem, err := memory.Get()
	before, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	total := float64(after.Total - before.Total)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	nodeStatus := model.Node{
		Name:                  CurrentInstallation.Name,
		LastSync:              time.Now().Unix(),
		Version:               command.Version(CurrentInstallation),
		NodeStatus:            command.NodeStatus(CurrentInstallation),
		PlotCount:             plotCount,
		PlotSize:              command.TotalPlotSize(CurrentInstallation),
		EstimatedTimeToWin:    command.EstimatedTimeToWin(CurrentInstallation),
		EstimatedNetworkSpace: command.EstimatedNetworkSpace(CurrentInstallation),
		FreeMemory:            int64(mem.Free),
		TotalMemory:           int64(mem.Total),
		CpuUsage:              int64(float64(after.System-before.System) / total * 100),
		DbSize:                0,
		Wallet:                "",
	}
	log.Println("Sending Node Status Update")
	jsonData, e := json.Marshal(nodeStatus)
	if e != nil {
		log.Println(e.Error())
	}
	req, err := http.NewRequest("POST", server+"/api/metric/node", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()
	log.Println("Status:" + response.Status)
}
