package main

import (
	"github.com/jasonlvhit/gocron"
	dbops "github.com/sooraj-sky/sky-meter/packages/dbops"
	skymeter "github.com/sooraj-sky/sky-meter/packages/httpserver"
	sentry "github.com/sooraj-sky/sky-meter/packages/logger"
	yamlops "github.com/sooraj-sky/sky-meter/packages/yamlops"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	log.Println("Launching sky-meter")
	sentry.SentryInit()
	dbconnect := os.Getenv("dbconnect")
	if opsgenieSecret == "" {
		log.Fatal("Please specify the opsgeniesecret as environment variable, e.g. sooraj@sky:~/go/src/sky-meter$ export dbconnect=host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
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
