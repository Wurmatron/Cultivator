package backend

import (
	"cultivator.wurmatron.io/backend/config"
	"cultivator.wurmatron.io/backend/pull"
	"cultivator.wurmatron.io/backend/routes"
	"cultivator.wurmatron.io/backend/storage"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

var DB *gorm.DB

func Start() {
	log.SetPrefix("[Backend]   > ")
	log.Println("Starting up as 'Backend'")
	DB = setupDBConnection()
	config.LoadOrSetupConfiguration(DB)
	router := mux.NewRouter()
	log.Println("Listening on '" + GetHost() + "' '" + "http://" + GetHost() + "/'")
	router = setupRoutes(*router)
	pull.CheckForTrackedUPSAndStartup()
	log.Fatal(http.ListenAndServe(GetHost(), router))
}

func GetHost() string {
	return config.Host + ":" + strconv.FormatInt(config.Port, 10)
}

func setupDBConnection() *gorm.DB {
	log.Printf("Connecting to db (%s) as user (%s) on (%s)", storage.Name, storage.User, storage.GetDBHost())
	DB, err := storage.GetConnection()
	if err != nil {
		log.Println("Failed to connect to db!")
		log.Fatal(err)
	}
	log.Println("Connection to DB has been established")
	return DB
}

func setupRoutes(router mux.Router) *mux.Router {
	// Information
	router.Handle("/api/metric/prometheus", promhttp.Handler())
	routes.AddMetricRoutes(&router)
	return &router
}
