package main

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/jasonlvhit/gocron"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	dbops "sky-meter/dbops"
	httpreponser "sky-meter/httpres"
	models "sky-meter/models"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to sky-meter")
}

func SelfStatusLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status OK")
}

func getStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	httpresdata, _ := httpreponser.GetHttpdata("http://apache.org")
	w.Write(httpresdata)
	return
}

func httpSyntheticCheck(endpoint string, time uint64) {
	gocron.Every(time).Second().Do(callEndpoint, endpoint)
	<-gocron.Start()
}

func callEndpoint(endpoint string) {
	httpresdata, _ := httpreponser.GetHttpdata(endpoint)
	log.Println(string(httpresdata))
}

func main() {
	sentenv := os.Getenv("sentry_dsn")
	if sentenv == "" {
		log.Fatal("Please specify the sentry_dsn as environment variable, e.g. env sentry_dsn=https://your-dentry-dsn.com go run server.go")
	}
	senterr := sentry.Init(sentry.ClientOptions{
		Dsn: sentenv,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if senterr != nil {
		log.Fatalf("sentry.Init: %s", senterr)
	}

	jsonFile, err := os.Open("input.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var endpoints models.JsonInput

	json.Unmarshal(byteValue, &endpoints)

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
	log.Println(urls)



	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Please specify the HTTP port as environment variable, e.g. env PORT=8081 go run http-server.go")
	}

	log.Println("listening on port", port)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/stats", getStats).Methods("GET")
	router.HandleFunc("/health", SelfStatusLink)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
