package routes

import (
	"cultivator.wurmatron.io/backend/config"
	"cultivator.wurmatron.io/backend/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/* Latest */

var NodeStatusCache map[string]model.Node
var HarvesterStatusCache map[string]model.Harvester
var PowerStatusCache map[string]model.Power

/* Temp for Historical */

var NodeHistory []model.Node
var HarvesterHistory []model.Harvester
var PowerHistory []model.Power

func init() {
	NodeStatusCache = make(map[string]model.Node)
	HarvesterStatusCache = make(map[string]model.Harvester)
	PowerStatusCache = make(map[string]model.Power)
}

func AddMetricRoutes(router *mux.Router, DB *gorm.DB) {
	// Node
	router.HandleFunc("/api/metric/node", UpdateNodeStatus).Methods("POST")
	router.HandleFunc("/api/metric/node/{name}", NodeGet).Methods("GET")
	router.HandleFunc("/api/metric/node", NodeGetAll).Methods("GET")
	// Harvester
	router.HandleFunc("/api/metric/harvester", UpdateHarvesterStatus).Methods("POST")
	router.HandleFunc("/api/metric/harvester/{id}", HarvesterGet).Methods("GET")
	router.HandleFunc("/api/metric/harvester", HarvesterGetAll).Methods("GET")
	router.HandleFunc("/api/metric/harvester", HarvesterGetAll).Methods("GET").Queries("name", "{name}")
	// Power
	router.HandleFunc("/api/metric/power/{serial}", PowerGet).Methods("GET")
	router.HandleFunc("/api/metric/power", PowerGetAll).Methods("GET")
	// Timeout / Cleanup
	go removeExpiredCache()
	// SQL History
	go handleSQLHistory(DB)
}

func removeExpiredCache() {
	for range time.Tick(time.Second * time.Duration(config.StatusTimeoutInterval)) {
		for key, val := range NodeStatusCache {
			if (val.LastSync + int64(config.StatusTimeoutInterval)) < time.Now().Unix() {
				log.Println("Node '" + key + "' has not sent a update in '" + strconv.FormatInt(config.StatusTimeoutInterval, 10) + "s' considering it expired / dead!")
				NodeHistory = append(NodeHistory, val)
				delete(NodeStatusCache, strings.ToLower(key))
				log.Println(key + " has been removed from the active node list")
			}
		}
		for key, val := range HarvesterStatusCache {
			if (val.LastSync + int64(config.StatusTimeoutInterval)) < time.Now().Unix() {
				log.Println("Harvester '" + val.ID + "' (" + key + ") has not sent a update in '" + strconv.FormatInt(config.StatusTimeoutInterval, 10) + "s' considering it expired / dead!")
				HarvesterHistory = append(HarvesterHistory, val)
				delete(HarvesterStatusCache, strings.ToLower(key))
				log.Println(key + " has been removed from the active harvester list")
			}
		}
		for key, val := range PowerStatusCache {
			if (val.Generated + int64(config.StatusTimeoutInterval)) < time.Now().Unix() {
				log.Println("Power '" + key + "' has not been responding in '" + strconv.FormatInt(config.StatusTimeoutInterval, 10) + "s' considering it dead / errored!")
				PowerHistory = append(PowerHistory, val)
				delete(PowerStatusCache, strings.ToLower(key))
				log.Println(key + " has been removed from the active Power list")
			}
		}
	}
}

func UpdateNodeStatus(w http.ResponseWriter, r *http.Request) {
	// Read json
	n := model.Node{}
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		log.Println("Failed to parse (" + err.Error() + ")")
		http.Error(w, "Failed to parse json object", http.StatusBadRequest)
	}
	// Process node update
	if isValidNodeStatusUpdate(n) {
		lower := strings.ToLower(n.Name)
		if hasNodeStatus(lower) {
			NodeHistory = append(NodeHistory, NodeStatusCache[lower])
		}
		NodeStatusCache[lower] = n
		log.Println("Node '" + n.Name + "' has been updated!")
	} else {
		http.Error(w, "Invalid node status update", http.StatusBadRequest)
		return
	}
	// Response
	w.WriteHeader(http.StatusAccepted)
}

func isValidNodeStatusUpdate(n model.Node) bool {
	return len(n.Name) > 0 && n.LastSync <= time.Now().Unix() && n.LastSync+(int64(config.StatusTimeoutInterval)*1000) > time.Now().Unix()
}

func hasNodeStatus(node string) bool {
	lower := strings.ToLower(node)
	if _, ok := NodeStatusCache[lower]; ok {
		return true
	}
	return false
}

func hasHarvesterStatus(node string) bool {
	lower := strings.ToLower(node)
	if _, ok := HarvesterStatusCache[lower]; ok {
		return true
	}
	return false
}

func NodeGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	if hasNodeStatus(name) {
		n := NodeStatusCache[strings.ToLower(name)]
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(n)
		if err != nil {
			log.Println("Failed to encode json for node '" + name + "'")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func NodeGetAll(w http.ResponseWriter, _ *http.Request) {
	arr := make([]model.Node, 0)
	for _, val := range NodeStatusCache {
		arr = append(arr, val)
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(arr)
	if err != nil {
		log.Println("Failed to encode json for nodes")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func HarvesterGetAll(w http.ResponseWriter, r *http.Request) {
	arr := make([]model.Harvester, 0)
	nameFilter := r.URL.Query().Get("name")
	if len(nameFilter) > 0 {
		for _, val := range HarvesterStatusCache {
			if strings.EqualFold(nameFilter, val.Blockchain) {
				arr = append(arr, val)
			}
		}
	} else {
		for _, val := range HarvesterStatusCache {
			arr = append(arr, val)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(arr)
	if err != nil {
		log.Println("Failed to encode json for harvesters")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func HarvesterGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if hasHarvesterStatus(id) {
		n := HarvesterStatusCache[strings.ToLower(id)]
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(n)
		if err != nil {
			log.Println("Failed to encode json for harvester '" + id + "'")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func UpdateHarvesterStatus(w http.ResponseWriter, r *http.Request) {
	// Read json
	h := model.Harvester{}
	if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
		log.Println("Failed to parse (" + err.Error() + ")")
		http.Error(w, "Failed to parse json object", http.StatusBadRequest)
	}
	// Process harvester update
	if isValidHarvesterStatusUpdate(h) {
		lower := strings.ToLower(h.Blockchain)
		if hasHarvesterStatus(lower) {
			HarvesterHistory = append(HarvesterHistory, HarvesterStatusCache[lower])
		}
		HarvesterStatusCache[lower] = h
		log.Println("Harvester '" + h.ID + "' has been updated!")
	} else {
		http.Error(w, "Invalid harvester status update", http.StatusBadRequest)
		return
	}
	// Response
	w.WriteHeader(http.StatusAccepted)
}

func isValidHarvesterStatusUpdate(n model.Harvester) bool {
	return len(n.Blockchain) > 0 && n.LastSync <= time.Now().Unix() && n.LastSync+(int64(config.StatusTimeoutInterval)*1000) > time.Now().Unix()
}

func hasUPSStatus(serial string) bool {
	lower := strings.ToLower(serial)
	if _, ok := PowerStatusCache[lower]; ok {
		return true
	}
	return false
}

func PowerGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	serial := params["serial"]
	if hasUPSStatus(serial) {
		n := PowerStatusCache[strings.ToLower(serial)]
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(n)
		if err != nil {
			log.Println("Failed to encode json for power '" + serial + "'")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func PowerGetAll(w http.ResponseWriter, r *http.Request) {
	arr := make([]model.Power, 0)
	for _, val := range PowerStatusCache {
		arr = append(arr, val)
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(arr)
	if err != nil {
		log.Println("Failed to encode json for power")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

/*
	Metrics History
*/

func handleSQLHistory(DB *gorm.DB) {
	for range time.Tick(time.Second * time.Duration(config.SqlStorageResolution)) {
		log.Println("Computing Avgs for incremental")
		go handleNodeHistory(DB)
		go handleHarvesterHistory(DB)
		go handlePowerHistory(DB)
	}
}

func handleNodeHistory(DB *gorm.DB) {
	if len(NodeHistory) > 0 {
		// Compute Individual Node Avg
		totalHistoryPerNode := int64(float64(config.StatusTimeoutInterval) / 4) // TODO Read Node Sync Config
		avgNodeHistory := map[string]model.ComputedNodeAvg{}
		for _, node := range NodeHistory {
			if _, ok := avgNodeHistory[strings.ToLower(node.Name)]; !ok {
				avgNodeHistory[strings.ToLower(node.Name)] = computeNodeHistory(node.Name, totalHistoryPerNode)
			}
		}
		// Compute Node Avg
		nodesAvg := model.ComputedNodeAvg{}
		nodesAvg.Name = "*"
		for _, nodeAvg := range avgNodeHistory {
			nodesAvg.CpuUsage = nodeAvg.CpuUsage + nodeAvg.CpuUsage
			nodesAvg.FreeMemory = nodeAvg.FreeMemory + nodeAvg.FreeMemory
			nodesAvg.Uptime = nodeAvg.Uptime + nodeAvg.Uptime
			nodesAvg.PlotCount = nodeAvg.PlotCount + nodeAvg.PlotCount
			nodesAvg.DbSize = nodeAvg.DbSize + nodeAvg.DbSize
		}
		nodesAvg.CpuUsage = nodesAvg.CpuUsage / int64(len(avgNodeHistory))
		nodesAvg.FreeMemory = nodesAvg.FreeMemory / int64(len(avgNodeHistory))
		nodesAvg.Uptime = nodesAvg.Uptime / float64(len(avgNodeHistory))
		nodesAvg.PlotCount = nodesAvg.PlotCount / int64(len(avgNodeHistory))
		nodesAvg.DbSize = nodesAvg.DbSize / int64(len(avgNodeHistory))
		// Send SQL for Storage
		sqlMetricsData := make([]model.Metrics, 0)
		for _, avg := range avgNodeHistory {
			sqlMetricsData = append(sqlMetricsData, convertToMetrics(avg, "increment", "node")...)
		}
		sqlMetricsData = append(sqlMetricsData, convertToMetrics(nodesAvg, "increment_avg", "node")...)
		DB.Table("metrics").CreateInBatches(sqlMetricsData, 100)
	} else {
		log.Println("No Node History Generated / Reported!")
	}
}

func convertToMetrics(avg model.ComputedNodeAvg, entryType string, style string) []model.Metrics {
	metrics := make([]model.Metrics, 0)
	metrics = append(metrics, createMetrics(avg.Name, style, entryType, "uptime", fmt.Sprintf("%.4f", avg.Uptime)))
	metrics = append(metrics, createMetrics(avg.Name, style, entryType, "cpu_usage", fmt.Sprintf("%d", avg.CpuUsage)))
	metrics = append(metrics, createMetrics(avg.Name, style, entryType, "plot_count", fmt.Sprintf("%d", avg.PlotCount)))
	metrics = append(metrics, createMetrics(avg.Name, style, entryType, "free_memory", fmt.Sprintf("%d", avg.FreeMemory)))
	metrics = append(metrics, createMetrics(avg.Name, style, entryType, "db_size", fmt.Sprintf("%d", avg.DbSize)))
	return metrics
}

func createMetrics(blockchain string, style string, entryType string, t string, val string) model.Metrics {
	return model.Metrics{
		Blockchain: strings.ToLower(blockchain),
		Style:      style,
		EntryType:  entryType,
		Type:       t,
		Value:      val,
		Timestamp:  fmt.Sprintf("%d", time.Now().Unix()),
	}
}

func countNodeEntries(name string) int64 {
	count := int64(0)
	for _, node := range NodeHistory {
		if strings.EqualFold(node.Name, name) {
			count++
		}
	}
	return count
}

func computeNodeHistory(name string, maxNodeEntries int64) model.ComputedNodeAvg {
	avg := model.ComputedNodeAvg{}
	nodeEntries := countNodeEntries(name)
	// Add Together
	avg.Name = name
	for _, node := range NodeHistory {
		if strings.EqualFold(node.Name, name) {
			avg.PlotCount = avg.PlotCount + node.PlotCount
			avg.FreeMemory = avg.FreeMemory + node.FreeMemory
			avg.CpuUsage = avg.PlotCount + node.CpuUsage
			avg.DbSize = avg.DbSize + node.DbSize
		}
	}
	// Div by count
	avg.PlotCount = avg.PlotCount / nodeEntries
	avg.FreeMemory = avg.FreeMemory / nodeEntries
	avg.CpuUsage = avg.CpuUsage / nodeEntries
	avg.DbSize = avg.DbSize / nodeEntries
	// Compute Uptime
	avg.Uptime = (float64(nodeEntries) / float64(maxNodeEntries)) * 100
	return avg
}

func handleHarvesterHistory(DB *gorm.DB) {
	if len(HarvesterHistory) > 0 {
		// Compute Individual Avg
		avgHarvesterHistory := map[string]model.ComputedHarvesterAvg{}
		for _, harvester := range HarvesterHistory {
			if _, ok := avgHarvesterHistory[strings.ToLower(harvester.ID)]; !ok {
				avgHarvesterHistory[harvester.ID] = computeHarvesterHistory(harvester.ID)
			}
		}
		// Compute Avg (Overall)
		harvesterAvg := model.ComputedHarvesterAvg{}
		harvesterAvg.ID = "*"
		harvesterAvg.Blockchain = "*"
		allDrives := []model.ComputedDriveAvg{}
		for _, harvester := range avgHarvesterHistory {
			harvesterAvg.PlotHarvestTime = harvesterAvg.PlotHarvestTime + harvester.PlotHarvestTime
			harvesterAvg.PlotCount = harvesterAvg.PlotCount + harvester.PlotCount
			// Drives
			allDrives = append(allDrives, harvester.Drives...)
		}
		avgDrive := model.ComputedDriveAvg{}
		avgDrive.Serial = "*"
		for _, drive := range allDrives {
			avgDrive.PlotCount = avgDrive.PlotCount + drive.PlotCount
			avgDrive.FreeSpace = avgDrive.FreeSpace + drive.FreeSpace
			avgDrive.TotalSpace = avgDrive.TotalSpace + drive.TotalSpace
			// Smart
			avgDrive.Smart.MaxTemp = avgDrive.Smart.MaxTemp + drive.Smart.MaxTemp
			avgDrive.Smart.MinTemp = avgDrive.Smart.MinTemp + drive.Smart.MinTemp
			avgDrive.Smart.CycleCount = avgDrive.Smart.CycleCount + drive.Smart.CycleCount
			avgDrive.Smart.CurrentTemp = avgDrive.Smart.CurrentTemp + drive.Smart.CurrentTemp
			avgDrive.Smart.PowerOnHours = avgDrive.Smart.PowerOnHours + drive.Smart.PowerOnHours
		}
		avgDrive.PlotCount = avgDrive.PlotCount / int64(len(allDrives))
		avgDrive.FreeSpace = avgDrive.FreeSpace / int64(len(allDrives))
		avgDrive.TotalSpace = avgDrive.TotalSpace / int64(len(allDrives))
		// Smart
		avgDrive.Smart.MaxTemp = avgDrive.Smart.MaxTemp / float64(len(allDrives))
		avgDrive.Smart.MinTemp = avgDrive.Smart.MinTemp / float64(len(allDrives))
		avgDrive.Smart.CycleCount = avgDrive.Smart.CycleCount / int64(len(allDrives))
		avgDrive.Smart.CurrentTemp = avgDrive.Smart.CurrentTemp / float64(len(allDrives))
		avgDrive.Smart.PowerOnHours = avgDrive.Smart.PowerOnHours / int64(len(allDrives))
		harvesterAvg.Drives = append(harvesterAvg.Drives, avgDrive)
		// Send SQL for Storage
		sqlMetricsData := make([]model.Metrics, 0)
		for _, avg := range avgHarvesterHistory {
			sqlMetricsData = append(sqlMetricsData, convertToMetricsHarvester(avg, "increment", "harvester")...)
		}
		sqlMetricsData = append(sqlMetricsData, convertToMetricsHarvester(harvesterAvg, "increment_avg", "harvester")...)
		DB.Table("metrics").CreateInBatches(sqlMetricsData, 100)
	} else {
		log.Println("No Harvester History Generated / Created!")
	}
}

func convertToMetricsHarvester(avg model.ComputedHarvesterAvg, entryType string, style string) []model.Metrics {
	metrics := make([]model.Metrics, 0)
	metrics = append(metrics, createMetrics(avg.ID, style, entryType, "plot_count", fmt.Sprintf("%d", avg.PlotCount)))
	metrics = append(metrics, createMetrics(avg.ID, style, entryType, "harvest_time", fmt.Sprintf("%.8f", avg.PlotHarvestTime)))
	metrics = append(metrics, createMetrics(avg.ID, style, entryType, "blockchain", fmt.Sprintf("%s", avg.Blockchain)))
	for _, drive := range avg.Drives {
		metrics = append(metrics, createMetrics(avg.ID, style, entryType, "drive_plot_count_"+drive.Serial, fmt.Sprintf("%d", drive.PlotCount)))
		metrics = append(metrics, createMetrics(avg.ID, style, entryType, "drive_total_space_"+drive.Serial, fmt.Sprintf("%d", drive.TotalSpace)))
		metrics = append(metrics, createMetrics(avg.ID, style, entryType, "drive_free_space_"+drive.Serial, fmt.Sprintf("%d", drive.FreeSpace)))
		metrics = append(metrics, createMetrics(avg.ID, style, entryType, "drive_cycle_count_"+drive.Serial, fmt.Sprintf("%d", drive.Smart.CycleCount)))
		metrics = append(metrics, createMetrics(avg.ID, style, entryType, "drive_max_temp_"+drive.Serial, fmt.Sprintf("%.4f", drive.Smart.MaxTemp)))
		metrics = append(metrics, createMetrics(avg.ID, style, entryType, "drive_min_temp_"+drive.Serial, fmt.Sprintf("%.4f", drive.Smart.MinTemp)))
		metrics = append(metrics, createMetrics(avg.ID, style, entryType, "drive_current_temp_"+drive.Serial, fmt.Sprintf("%.4f", drive.Smart.CurrentTemp)))
		metrics = append(metrics, createMetrics(avg.ID, style, entryType, "drive_power_on_hours_"+drive.Serial, fmt.Sprintf("%d", drive.Smart.PowerOnHours)))
	}
	return metrics
}

func computeHarvesterHistory(ID string) model.ComputedHarvesterAvg {
	avg := model.ComputedHarvesterAvg{}
	avg.ID = ID
	allDrives := []model.Drive{}
	for _, harvester := range HarvesterHistory {
		if strings.EqualFold(harvester.ID, ID) {
			avg.Blockchain = harvester.Blockchain
			avg.PlotCount = avg.PlotCount + harvester.PlotCount
			avg.PlotHarvestTime = avg.PlotHarvestTime + harvester.PlotHarvestTime
			allDrives = append(allDrives, harvester.Drives...)
		}
	}
	for _, drive := range allDrives {
		avg.Drives = append(avg.Drives, computeDriveAvg(drive.Serial, allDrives))
	}
	// Div by count
	harvesterEntries := countHarvesterEntries(ID)
	avg.PlotCount = avg.PlotCount / harvesterEntries
	avg.PlotHarvestTime = avg.PlotHarvestTime / float64(harvesterEntries)
	return avg
}

func countHarvesterEntries(ID string) int64 {
	count := 0
	for _, x := range HarvesterHistory {
		if strings.EqualFold(ID, x.ID) {
			count++
		}
	}
	return int64(count)
}

func computeDriveAvg(serial string, drives []model.Drive) model.ComputedDriveAvg {
	driveAvg := model.ComputedDriveAvg{}
	driveAvg.Serial = serial
	count := 0
	for _, drive := range drives {
		if strings.EqualFold(drive.Serial, serial) {
			count++
			driveAvg.PlotCount = driveAvg.PlotCount + drive.PlotCount
			driveAvg.TotalSpace = driveAvg.TotalSpace + drive.TotalSpace
			driveAvg.FreeSpace = driveAvg.FreeSpace + drive.FreeSpace
			// Smart
			driveAvg.Smart.CycleCount = driveAvg.Smart.CycleCount + drive.Smart.CycleCount
			driveAvg.Smart.MaxTemp = driveAvg.Smart.MaxTemp + drive.Smart.MaxTemp
			driveAvg.Smart.CurrentTemp = driveAvg.Smart.CurrentTemp + drive.Smart.CurrentTemp
			driveAvg.Smart.PowerOnHours = driveAvg.Smart.PowerOnHours + drive.Smart.PowerOnHours
		}
	}
	driveAvg.PlotCount = driveAvg.PlotCount / int64(count)
	driveAvg.TotalSpace = driveAvg.TotalSpace / int64(count)
	driveAvg.FreeSpace = driveAvg.FreeSpace / int64(count)
	driveAvg.Smart.CycleCount = driveAvg.Smart.CycleCount / int64(count)
	driveAvg.Smart.MaxTemp = driveAvg.Smart.MaxTemp / float64(count)
	driveAvg.Smart.CurrentTemp = driveAvg.Smart.CurrentTemp / float64(count)
	driveAvg.Smart.PowerOnHours = driveAvg.Smart.PowerOnHours / int64(count)
	return driveAvg
}

func handlePowerHistory(DB *gorm.DB) {

}
