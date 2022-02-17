package command

import (
	"bytes"
	"cultivator.wurmatron.io/backend/model"
	"log"
	"os"
	"os/exec"
	"strings"
)

func RunCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Println("Failed to run command '" + name + " " + strings.Join(arg, " ") + "'")
		log.Println(err.Error())
	}
}

func RunNodeCommand(blockchain model.BlockchainInstallation, command string) []string {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("sh", "-c", "cd "+blockchain.InstallDir+" && . ./activate && "+strings.ToLower(blockchain.Name)+" "+command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	_ = cmd.Run()
	arr := make([]string, 0)
	arr = append(arr, strings.Split(stdout.String(), "\n")...)
	return arr
}

func Exists(d string) bool {
	dir, err := os.Stat(d)
	if os.IsNotExist(err) {
		return false
	}
	return dir.IsDir()
}

func StartNode(chain model.BlockchainInstallation) {
	RunNodeCommand(chain, "start node")
}

func StopNode(chain model.BlockchainInstallation) {
	RunNodeCommand(chain, "stop node")
}

func StartHarvester(chain model.BlockchainInstallation) {
	RunNodeCommand(chain, "start harvester")
}

func infoCollector(chain model.BlockchainInstallation, command string, prefix string) string {
	status := RunNodeCommand(chain, command)
	for _, s := range status {
		if strings.HasPrefix(s, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(s, prefix))
		}
	}
	return "Failed to find '" + prefix + "' with " + command
}

func NodeStatus(chain model.BlockchainInstallation) string {
	return infoCollector(chain, "show -s", "Current Blockchain Status:")
}

func EstimatedNetworkSpace(chain model.BlockchainInstallation) string {
	return infoCollector(chain, "show -s", "Estimated network space:")
}

func TotalPlots(chain model.BlockchainInstallation) string {
	return infoCollector(chain, "farm summary", "Plot count for all harvesters:")
}

func TotalPlotSize(chain model.BlockchainInstallation) string {
	return infoCollector(chain, "farm summary", "Total size of plots:")
}

func EstimatedTimeToWin(chain model.BlockchainInstallation) string {
	return infoCollector(chain, "farm summary", "Expected time to win:")
}

func Version(chain model.BlockchainInstallation) string {
	return infoCollector(chain, "version", "")
}
