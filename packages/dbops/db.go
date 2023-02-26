package dbops

import (
	"encoding/json"
	"log"
	models "sky-meter/models"
	skyalerts "sky-meter/packages/alerts"
	httpreponser "sky-meter/packages/httpres"
	"time"

	"gorm.io/gorm"
)

// Define an interface named error to avoid name collision with the built-in error interface
type error interface {
	Error() string
}

// Perform initial migrations for the database tables
func InitialMigration(db *gorm.DB) {
	db.AutoMigrate(&models.AllEndpoints{})
	db.AutoMigrate(&models.HttpOutput{})
	db.AutoMigrate(&models.OpsgenieAlertData{})
	db.AutoMigrate(&models.AlertGroups{})
}

// Insert URLs and associated data into the database
func InsertUrlsToDb(db *gorm.DB, endpoints models.UserInput) {
	var urlCheck models.AllEndpoints
	var urlsId models.AllEndpoints
	var Groups models.AlertGroups

	// Loop through each domain and check if it already exists in the database; if it does not, create it
	for i := 0; i < len(endpoints.Domains); i++ {
		db.Where("URL=?", endpoints.Domains[i].Name).Find(&urlCheck)
		if urlCheck.CreatedAt == 0 && urlCheck.URL != endpoints.Domains[i].Name {
			db.Create(&models.AllEndpoints{URL: endpoints.Domains[i].Name, Timeout: endpoints.Domains[i].Timeout, SkipSsl: endpoints.Domains[i].SkipSsl, Frequency: endpoints.Domains[i].Frequency, Group: endpoints.Domains[i].Group, Active: true})
		}
		urlCheck = urlsId
	}

	// Loop through each group and its associated emails, and create them if they do not exist in the database
	for i := range endpoints.Groups {
		db.Where("NAME=?", endpoints.Groups[i].Name).Find(&Groups)
		for k := range endpoints.Groups[i].Emails {
			if Groups.CreatedAt == 0 && Groups.Email != endpoints.Groups[i].Emails[k] {
				db.Create(&models.AlertGroups{Name: endpoints.Groups[i].Name, Email: endpoints.Groups[i].Emails[k]})
			}
		}

	}

}

// Check the frequency of each URL in the database, and send alerts if necessary
func GetUrlFrequency(db *gorm.DB, endpoints models.UserInput) {
	var urlsToCheck []models.AllEndpoints
	var urlsId []models.AllEndpoints
	var alertStatus models.OpsgenieAlertData
	var GroupsEmailIds []models.AlertGroups
	var Empty []models.AlertGroups

	// Find all URLs in the database
	db.Find(&urlsToCheck)

	// Loop through each URL
	for i := 0; i < len(urlsToCheck); i++ {

		// Retrieve the URL by ID
		db.First(&urlsId, urlsToCheck[i].ID)

		// Check if the URL is active
		if urlsToCheck[i].Active {

			// If the URL has not yet been checked, update the "next run" field and check it
			if urlsToCheck[i].NextRun == 0 {
				db.Model(&urlsId).Where("id = ?", urlsToCheck[i].ID).Update("next_run", urlsToCheck[i].Frequency)

				// Call the URL and retrieve the HTTP output and status code
				httpOutput, HttpStatusCode, err := httpreponser.CallEndpoint(urlsToCheck[i].URL, urlsToCheck[i].Timeout, urlsToCheck[i].SkipSsl)

				// If an error occurred, send alerts
				if err != nil {

					// Retrieve any existing alerts for the URL
					db.First(&alertStatus, "url = ?", urlsToCheck[i].URL)

					// Find alert groups for the URL
					db.Where("Name=?", urlsToCheck[i].Group).Find(&GroupsEmailIds)
					var emailIds []string
					for _, group := range GroupsEmailIds {
						emailIds = append(emailIds, group.Email)
					}

					// Check if there is already an existing alert for the URL
					if alertStatus.URL == urlsToCheck[i].URL {
						dt := time.Now()

						// Send Opsgenie Notification
						if endpoints.Opegenie.Enabled {
							AlertStatus := skyalerts.CheckAlertStatus(alertStatus.RequestId)
							if (AlertStatus == "closed") || (alertStatus.Error != err.Error()) {
								alertReqId := skyalerts.OpsgenieCreateAlert(urlsToCheck[i].URL, err, urlsToCheck[i].Group)
								db.Model(&alertStatus).Where("url = ?", urlsToCheck[i].URL).Update("request_id", alertReqId)
							} else {
								alertReqId := "Opegenie-disabled"
								db.Model(&alertStatus).Where("url = ?", urlsToCheck[i].URL).Update("request_id", alertReqId)
							}

							d := models.SmtpErr{urlsToCheck[i].URL, "webiste is Down", dt, err.Error(), emailIds}

							// Send Email Notification
							if endpoints.Email.Enabled {
								skyalerts.SendMail(d)
							}

						}

						// Reset alert groups
						GroupsEmailIds = Empty
					} else {

						// If there is no existing alert for the URL, create a new one
						dts := time.Now()
						db.Create(&models.HttpOutput{OutputData: httpOutput, URL: urlsToCheck[i].URL, StatusCode: HttpStatusCode, Error: err.Error()})
						d := models.SmtpErr{urlsToCheck[i].URL, "webiste is Down", dts, err.Error(), emailIds}
						if endpoints.Opegenie.Enabled {
							alertReqId := skyalerts.OpsgenieCreateAlert(urlsToCheck[i].URL, err, urlsToCheck[i].Group)
							db.Create(&models.OpsgenieAlertData{URL: urlsToCheck[i].URL, RequestId: alertReqId, Error: err.Error(), Active: true})
							db.Create(&models.HttpOutput{OutputData: httpOutput, URL: urlsToCheck[i].URL, StatusCode: HttpStatusCode, Error: err.Error()})
						} else {
							alertReqId := "Opegenie-disabled"
							db.Create(&models.OpsgenieAlertData{URL: urlsToCheck[i].URL, RequestId: alertReqId, Error: err.Error(), Active: true})
							db.Create(&models.HttpOutput{OutputData: httpOutput, URL: urlsToCheck[i].URL, StatusCode: HttpStatusCode, Error: err.Error()})

						}

						// Send Email Notification
						if endpoints.Email.Enabled {
							skyalerts.SendMail(d)
						}
					}

				} else {
					var byteHttpOutput models.Debug
					json.Unmarshal(httpOutput, &byteHttpOutput)
					db.Create(&models.HttpOutput{OutputData: httpOutput, URL: urlsToCheck[i].URL, StatusCode: HttpStatusCode})
				}

			} else {
				db.Model(&urlsId).Where("id = ?", urlsToCheck[i].ID).Update("next_run", urlsToCheck[i].NextRun-1)
			}
		}
	}
}

// This function removes old entries from the database based on the input provided.
// It takes a database object and user input object as parameters.
func RemoveOldEntry(db *gorm.DB, endpoints models.UserInput) {

	// Retrieve all endpoints from the database.
	var urlCheck []models.AllEndpoints
	db.Find(&urlCheck)
	var urlsId models.AllEndpoints
	// Iterate through all endpoints from the database.
	for i := 0; i < len(urlCheck); i++ {

		// Initialize a variable to keep track of how many endpoints are found.
		var findCount int

		// Iterate through all domains in the user input object.
		for m := range endpoints.Domains {

			// Check if the domain from user input is present in the database.
			if endpoints.Domains[m].Name == urlCheck[i].URL {
				log.Println("found URL", urlCheck[i].URL)

				// If the endpoint is inactive, update its status to active.
				if urlCheck[i].Active != true {
					db.Model(&urlsId).Where("url = ?", urlCheck[i].URL).Update("Active", true)
				}
			} else {
				findCount = findCount + 1
			}
		}

		// If the endpoint is not found in the user input object, update its status to inactive.
		if findCount >= len(endpoints.Domains) {
			log.Println(urlCheck[i].URL, "Not found in input.json, removing the check")
			db.Model(&urlsId).Where("url = ?", urlCheck[i].URL).Update("Active", false)
		}
	}

}
