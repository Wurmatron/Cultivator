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
				log.Println("Harvester '" + key + "' has not sent a update in '" + strconv.FormatInt(config.StatusTimeoutInterval, 10) + "s' considering it expired / dead!")
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
			if strings.EqualFold(nameFilter, val.Name) {
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
		lower := strings.ToLower(h.Name)
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
	return len(n.Name) > 0 && n.LastSync <= time.Now().Unix() && n.LastSync+(int64(config.StatusTimeoutInterval)*1000) > time.Now().Unix()
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
		go handleUPSHistory(DB)
	}
}

func handleNodeHistory(DB *gorm.DB) {
	// Compute Individual Node Avg
	totalHistoryPerNode := int64(float64(config.StatusTimeoutInterval) / 4) // TODO Read Node Sync Config
	avgNodeHistory := map[string]model.ComputedNodeAvg{}
	for _, node := range NodeHistory {
		if _, ok := PowerStatusCache[strings.ToLower(node.Name)]; !ok {
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
		sqlMetricsData = append(sqlMetricsData, convertToMetrics(avg, "increment")...)
	}
	sqlMetricsData = append(sqlMetricsData, convertToMetrics(nodesAvg, "avg")...)
	DB.Table("metrics").CreateInBatches(sqlMetricsData, 100)
}

func convertToMetrics(avg model.ComputedNodeAvg, entryType string) []model.Metrics {
	metrics := make([]model.Metrics, 0)
	metrics = append(metrics, createMetrics(avg.Name, entryType, "uptime", fmt.Sprintf("%.4f", avg.Uptime)))
	metrics = append(metrics, createMetrics(avg.Name, entryType, "cpu_usage", fmt.Sprintf("%d", avg.CpuUsage)))
	metrics = append(metrics, createMetrics(avg.Name, entryType, "plot_count", fmt.Sprintf("%d", avg.PlotCount)))
	metrics = append(metrics, createMetrics(avg.Name, entryType, "free_memory", fmt.Sprintf("%d", avg.FreeMemory)))
	metrics = append(metrics, createMetrics(avg.Name, entryType, "db_size", fmt.Sprintf("%d", avg.DbSize)))
	return metrics
}

func createMetrics(blockchain string, entryType string, t string, val string) model.Metrics {
	return model.Metrics{
		Blockchain: strings.ToLower(blockchain),
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

}

func handleUPSHistory(DB *gorm.DB) {

}
