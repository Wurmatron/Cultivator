package routes

import (
	"cultivator.wurmatron.io/backend/config"
	"cultivator.wurmatron.io/backend/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var NodeStatusCache map[string]model.Node
var HarvesterStatusCache map[string]model.Harvester

func init() {
	NodeStatusCache = make(map[string]model.Node)
	HarvesterStatusCache = make(map[string]model.Harvester)
}

func AddMetricRoutes(router *mux.Router) {
	// Node
	router.HandleFunc("/api/metric/node", UpdateNodeStatus).Methods("POST")
	router.HandleFunc("/api/metric/node/{name}", NodeGet).Methods("GET")
	router.HandleFunc("/api/metric/node", NodeGetAll).Methods("GET")
	// Harvester
	router.HandleFunc("/api/metric/harvester", UpdateHarvesterStatus).Methods("POST")
	router.HandleFunc("/api/metric/harvester/{id}", HarvesterGet).Methods("GET")
	router.HandleFunc("/api/metric/harvester", HarvesterGetAll).Methods("GET")
	router.HandleFunc("/api/metric/harvester", HarvesterGetAll).Methods("GET").Queries("name", "{name}")
	// Timeout / Cleanup
	go removeExpiredCache()
}

func removeExpiredCache() {
	for range time.Tick(time.Second * time.Duration(config.StatusTimeoutInterval)) {
		for key, val := range NodeStatusCache {
			if (val.LastSync + int64(config.StatusTimeoutInterval)) < time.Now().Unix() {
				log.Println("Node '" + key + "' has not sent a update in '" + strconv.FormatInt(config.StatusTimeoutInterval, 10) + "s' considering it expired / dead!")
				delete(NodeStatusCache, strings.ToLower(key))
				log.Println(key + " has been removed from the active node list")
			}
		}
		for key, val := range HarvesterStatusCache {
			if (val.LastSync + int64(config.StatusTimeoutInterval)) < time.Now().Unix() {
				log.Println("Harvester '" + key + "' has not sent a update in '" + strconv.FormatInt(config.StatusTimeoutInterval, 10) + "s' considering it expired / dead!")
				delete(HarvesterStatusCache, strings.ToLower(key))
				log.Println(key + " has been removed from the active harvester list")
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
