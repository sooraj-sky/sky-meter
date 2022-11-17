package main

import (
	"log"
	dbops "sky-meter/packages/dbops"
	jsonops "sky-meter/packages/jsonops"
	sentry "sky-meter/packages/logger"

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

	}), &gorm.Config{ Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Println(err)
	}

	endpoints := jsonops.InputJson()

	dbops.InitialMigration(db)
	dbops.InsertUrlsToDb(db, endpoints)
	log.Println("Updated sky-meter targets")
	log.Println("Staring sky-meter Health Check")
	gocron.Every(1).Second().Do(dbops.GetUrlFrequency, db)
	<-gocron.Start()
	//	skymeter.InitServer()

}
