package main

import (
	"github.com/jasonlvhit/gocron"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	dbops "sky-meter/packages/dbops"
	skyenv "sky-meter/packages/env"
	sentry "sky-meter/packages/logger"
	yamlops "sky-meter/packages/yamlops"
)

func init() {

	// Initialize the environment variables
	skyenv.InitEnv()
}

func main() {
	log.Println("Launching sky-meter")

	// Initialize the Sentry logger
	sentry.SentryInit()

	// Get all the environment variables
	allEnv := skyenv.GetEnv()
	dbconnect := allEnv.DbUrl
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbconnect,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage

	}), &gorm.Config{})

	if err != nil {

		// Exit the program if there is an error connecting to the database
		log.Fatal(err)
	}

	// Read the YAML file to get the list of endpoints to monitor
	endpoints := yamlops.InputYml()
	dbops.InitialMigration(db)
	dbops.InsertUrlsToDb(db, endpoints)
	dbops.RemoveOldEntry(db, endpoints)
	log.Println("Updated sky-meter targets")
	log.Println("Staring sky-meter Health Check")
	gocron.Every(1).Second().Do(dbops.GetUrlFrequency, db, endpoints)
	<-gocron.Start() // Start the scheduler and block the main thread until it exits

}
