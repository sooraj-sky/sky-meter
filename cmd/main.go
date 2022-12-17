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
