package main

import (
	"github.com/gorilla/mux"
	"github.com/jasonlvhit/gocron"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	api "sky-meter/packages/api"
	dbops "sky-meter/packages/dbops"
	jsonops "sky-meter/packages/jsonops"
	sentry "sky-meter/packages/logger"
)

func main() {
	sentry.SentryInit()
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage

	}), &gorm.Config{})

	if err != nil {
		log.Println(err)
	}

	endpoints := jsonops.InputJson()
	dbops.InitialMigration(db)
	dbops.InsertUrlsToDb(db, endpoints)

	gocron.Every(1).Second().Do(dbops.GetUrlFrequency, db)
	<-gocron.Start()

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
