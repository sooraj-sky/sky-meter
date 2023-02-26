package main

import (
	"github.com/jasonlvhit/gocron"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	dbops "sky-meter/packages/dbops"
	skyenv "sky-meter/packages/env"
	skymeter "sky-meter/packages/httpserver"
	sentry "sky-meter/packages/logger"
	yamlops "sky-meter/packages/yamlops"
)

func init() {
	skyenv.InitEnv()
}

func main() {
	log.Println("Launching sky-meter")
	sentry.SentryInit()
	allEnv := skyenv.GetEnv()
	dbconnect := allEnv.DbUrl
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
