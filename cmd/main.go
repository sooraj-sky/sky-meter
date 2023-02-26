package main

import (
	"github.com/jasonlvhit/gocron"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	dbops "sky-meter/packages/dbops"
	skymeter "sky-meter/packages/httpserver"
	sentry "sky-meter/packages/logger"
	yamlops "sky-meter/packages/yamlops"
)

func main() {
	log.Println("Launching sky-meter")
	sentry.SentryInit()
	dbconnect := os.Getenv("dbconnect")
	opsgenieSecret := os.Getenv("opsgeniesecret")
	if opsgenieSecret == "" {
		log.Fatal("Please specify the opsgeniesecret as environment variable, e.g.  export dbconnect=host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	}

	if opsgenieSecret == "" {
		log.Fatal("Please specify the opsgeniesecret as environment variable, e.g.  export opsgeniesecret=<your-value-here>")
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbconnect,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage

	}), &gorm.Config{})

	if err != nil {
		log.Println(err)
	}
	endpoints := yamlops.InputYml()
	dbops.InitialMigration(db)
	dbops.InsertUrlsToDb(db, endpoints)
	dbops.RemoveOldEntry(db, endpoints)
	log.Println("Updated sky-meter targets")
	log.Println("Staring sky-meter Health Check")
	gocron.Every(1).Second().Do(dbops.GetUrlFrequency, db)
	<-gocron.Start()
	skymeter.InitServer()

}
