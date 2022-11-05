package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jasonlvhit/gocron"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	dbops "sky-meter/packages/dbops"
	httpreponser "sky-meter/packages/httpres"
	models "sky-meter/models"
	api "sky-meter/packages/api"
	sentry "sky-meter/packages/logger"
	jsonops "sky-meter/packages/jsonops"
)

func httpSyntheticCheck(endpoint string, time uint64) {
	gocron.Every(time).Second().Do(callEndpoint, endpoint)
	<-gocron.Start()
}

func callEndpoint(endpoint string) {
	httpresdata, _ := httpreponser.GetHttpdata(endpoint)
	log.Println(string(httpresdata))
}

func main() {
	sentry.SentryInit()
    endpoints := jsonops.InputJson()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Println(err)
	}

	dbops.InitialMigration(db)
	dbops.InsertUrlsToDb(db, endpoints)
	urls := dbops.GetUrlFrequency(db)

	if urls != nil {
		url, _ := urls.(models.AllEndpoints)
		fmt.Println(url.URL)
		callEndpoint(url.URL)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Please specify the HTTP port as environment variable, e.g. env PORT=8081 go run http-server.go")
	}

	log.Println("listening on port", port)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", api.HomeLink)
	router.HandleFunc("/health", api.SelfStatusLink)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
