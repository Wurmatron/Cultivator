package pull

import (
	"cultivator.wurmatron.io/backend/config"
	"cultivator.wurmatron.io/backend/model"
	"cultivator.wurmatron.io/backend/routes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func CheckForTrackedUPSAndStartup() {
	if len(config.PowerPullIP) > 0 && len(config.PowerPullIP[0]) > 0 {
		go updatePowerDataPeriodically()
	}
}

func updatePowerDataPeriodically() {
	UpdatePowerData()
	for range time.Tick(time.Second * time.Duration(config.PowerPollFrequency)) {
		UpdatePowerData()
	}
}

func UpdatePowerData() {
	for _, ip := range config.PowerPullIP {
		res, err := http.Get(ip)
		if err != nil {
			log.Println("Failed to pull Power from '" + ip + "'")
		}
		jsonData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("Failed to read json data from Power Rest API! (" + ip + ")")
		}
		ups := model.Power{}
		err = json.Unmarshal(jsonData, &ups)
		if err != nil {
			log.Println("Failed to parse json data from Power Rest API! (" + ip + ")")
		}
		if hasPowerStatus(ups.Serial) {
			routes.PowerHistory = append(routes.PowerHistory, routes.PowerStatusCache[ups.Serial])
		}
		routes.PowerStatusCache[ups.Serial] = ups
	}
}

func hasPowerStatus(serial string) bool {
	lower := strings.ToLower(serial)
	if _, ok := routes.PowerStatusCache[lower]; ok {
		return true
	}
	return false
}
