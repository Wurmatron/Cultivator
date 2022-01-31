package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type UPS struct {
	MaxPowerDraw          int64   `json:"maxPowerDraw"`
	CurrentLoad           int64   `json:"currentLoad"`
	InputVoltage          float64 `json:"inputVoltage"`
	OutputVoltage         float64 `json:"outputVoltage"`
	InputVoltageNominal   float64 `json:"inputVoltageNominal"`
	Wattage               float64 `json:"wattage"`
	Model                 string  `json:"model"`
	Serial                string  `json:"serial"`
	BatteryVoltage        float64 `json:"batteryVoltage"`
	BatteryVoltageNominal float64 `json:"BatteryVoltageNominal"`
	BatteryType           string  `json:"batteryType"`
	BatteryCharge         int64   `json:"batteryCharge"`
	BatteryRuntime        int64   `json:"batteryRuntime"`
	Generated             int64   `json:"generated"`
}

const (
	ups_name = "power"
	port     = "8090"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/power/metrics", createJson)
	log.Println("Starting UPS API on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func createJson(w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold(r.Method, "GET") {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(convertToJson(strings.Split(getRawData(), "\n"))))
	}
}

func getRawData() string {
	out, err := exec.Command("upsc", ups_name).Output()
	if err != nil {
		log.Println("nut has not been configured")
		log.Fatal(err)
	}
	return string(out)
}

func convertToJson(raw []string) string {
	ups := UPS{
		MaxPowerDraw:          getDataInt64(raw, "power.realpower.nominal"),
		CurrentLoad:           getDataInt64(raw, "power.load"),
		InputVoltage:          getDataFloat64(raw, "input.voltage"),
		OutputVoltage:         getDataFloat64(raw, "output.voltage"),
		InputVoltageNominal:   getDataFloat64(raw, "input.voltage.nominal"),
		Wattage:               (float64(getDataInt64(raw, "power.load")) / 100.0) * getDataFloat64(raw, "power.realpower.nominal"),
		Model:                 getDataStr(raw, "power.model"),
		Serial:                getDataStr(raw, "power.serial"),
		BatteryVoltage:        getDataFloat64(raw, "battery.voltage"),
		BatteryVoltageNominal: getDataFloat64(raw, "battery.voltage.nominal"),
		BatteryType:           getDataStr(raw, "battery.type"),
		BatteryCharge:         getDataInt64(raw, "battery.charge"),
		BatteryRuntime:        getDataInt64(raw, "battery.runtime"),
		Generated:             time.Now().Unix(),
	}
	json, err := json.Marshal(ups)
	if err != nil {
		log.Println("Failed to convert raw data into json")
	}
	return string(json)
}

func getDataInt64(raw []string, str string) int64 {
	inr, err := strconv.ParseInt(getDataStr(raw, str), 10, 64)
	if err != nil {
		log.Println(err)
		return -1
	}
	return inr
}

func getDataFloat64(raw []string, str string) float64 {
	inr, err := strconv.ParseFloat(getDataStr(raw, str), 64)
	if err != nil {
		log.Println(err)
		return -1
	}
	return inr
}

func getDataStr(raw []string, str string) string {
	for _, s := range raw {
		if strings.HasPrefix(s, str) {
			s = strings.TrimPrefix(s, str)
			s = strings.Trim(s, ":")
			s = strings.Trim(s, " ")
			return s
		}
	}
	return ""
}
