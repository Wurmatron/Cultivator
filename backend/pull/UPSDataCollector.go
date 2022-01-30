package pull

import (
	"cultivator.wurmatron.io/backend/config"
	"cultivator.wurmatron.io/backend/model"
	"cultivator.wurmatron.io/backend/routes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func CheckForTrackedUPSAndStartup() {
	if len(config.UpsPullLocationIP) > 0 {
		go updateUPSDataPeriodically()
	}
}

func updateUPSDataPeriodically() {
	updateUPSData()
	for range time.Tick(time.Second * time.Duration(config.UPSPollFrequency)) {
		updateUPSData()
	}
}

func updateUPSData() {
	for _, ip := range config.UpsPullLocationIP {
		res, err := http.Get(ip)
		if err != nil {
			log.Println("Failed to pull UPS from '" + ip + "'")
		}
		jsonData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("Failed to read json data from UPS Rest API! (" + ip + ")")
		}
		ups := model.UPS{}
		err = json.Unmarshal(jsonData, &ups)
		if err != nil {
			log.Println("Failed to parse json data from UPS Rest API! (" + ip + ")")
		}
		routes.UPSStatusCache[ups.Serial] = ups
	}
}
