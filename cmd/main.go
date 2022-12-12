package main

import (
	"log"
	dbops "sky-meter/packages/dbops"
	skymeter "sky-meter/packages/httpserver"
	sentry "sky-meter/packages/logger"
	yamlops "sky-meter/packages/yamlops"

	"github.com/jasonlvhit/gocron"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("Launching sky-meter")
	sentry.SentryInit()
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable",
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
