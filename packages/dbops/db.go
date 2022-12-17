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

type error interface {
	Error() string
}

func InitialMigration(db *gorm.DB) {
	db.AutoMigrate(&models.AllEndpoints{})
	db.AutoMigrate(&models.HttpOutput{})
	db.AutoMigrate(&models.OpsgenieAlertData{})
	db.AutoMigrate(&models.AlertGroups{})
}

func InsertUrlsToDb(db *gorm.DB, endpoints models.UserInput) {
	var urlCheck models.AllEndpoints
	var urlsId models.AllEndpoints
	var Groups models.AlertGroups
	for i := 0; i < len(endpoints.Domains); i++ {
		db.Where("URL=?", endpoints.Domains[i].Name).Find(&urlCheck)
		if urlCheck.CreatedAt == 0 && urlCheck.URL != endpoints.Domains[i].Name {
			db.Create(&models.AllEndpoints{URL: endpoints.Domains[i].Name, Timeout: endpoints.Domains[i].Timeout, SkipSsl: endpoints.Domains[i].SkipSsl, Frequency: endpoints.Domains[i].Frequency, Group: endpoints.Domains[i].Group, Active: true})
		}
		urlCheck = urlsId
	}

	for i := range endpoints.Groups {
		db.Where("NAME=?", endpoints.Groups[i].Name).Find(&Groups)
		for k := range endpoints.Groups[i].Emails {
			if Groups.CreatedAt == 0 && Groups.Email != endpoints.Groups[1].Emails[k] {
				db.Create(&models.AlertGroups{Name: endpoints.Groups[i].Name, Email: endpoints.Groups[1].Emails[k]})
			}
		}

	}

}

func GetUrlFrequency(db *gorm.DB) {
	var urlsToCheck []models.AllEndpoints
	var urlsId []models.AllEndpoints
	var alertStatus models.OpsgenieAlertData
	var GroupsEmailIds models.AlertGroups
	var Empty models.AlertGroups

	db.Find(&urlsToCheck)
	for i := 0; i < len(urlsToCheck); i++ {
		db.First(&urlsId, urlsToCheck[i].ID)
		if urlsToCheck[i].Active {
			if urlsToCheck[i].NextRun == 0 {
				db.Model(&urlsId).Where("id = ?", urlsToCheck[i].ID).Update("next_run", urlsToCheck[i].Frequency)
				httpOutput, HttpStatusCode, err := httpreponser.CallEndpoint(urlsToCheck[i].URL, urlsToCheck[i].Timeout, urlsToCheck[i].SkipSsl)
				if err != nil {
					db.First(&alertStatus, "url = ?", urlsToCheck[i].URL)
					db.Where("Name=?",  urlsToCheck[i].Group).Find(&GroupsEmailIds)

					log.Println("woeeeeeeeeeee", GroupsEmailIds)

					if alertStatus.URL == urlsToCheck[i].URL {
						dt := time.Now()
						AlertStatus := skyalerts.CheckAlertStatus(alertStatus.RequestId)
						if (AlertStatus == "closed") || (alertStatus.Error != err.Error()) {
								alertReqId := skyalerts.OpsgenieCreateAlert(urlsToCheck[i].URL, err, urlsToCheck[i].Group)
								db.Model(&alertStatus).Where("url = ?", urlsToCheck[i].URL).Update("request_id", alertReqId)
							d := models.SmtpErr{urlsToCheck[i].URL, "webiste is Down", dt, err.Error(), []string{GroupsEmailIds.Email}}


							skyalerts.SendMail(d)
						}
                        GroupsEmailIds = Empty
					} else {

						dts := time.Now()
							db.Create(&models.HttpOutput{OutputData: httpOutput, URL: urlsToCheck[i].URL, StatusCode: HttpStatusCode, Error: err.Error()})
						d := models.SmtpErr{urlsToCheck[i].URL, "webiste is Down", dts, err.Error(), []string{GroupsEmailIds.Email}}

							alertReqId := skyalerts.OpsgenieCreateAlert(urlsToCheck[i].URL, err, urlsToCheck[i].Group)
							db.Create(&models.OpsgenieAlertData{URL: urlsToCheck[i].URL, RequestId: alertReqId, Error: err.Error(), Active: true})
							db.Create(&models.HttpOutput{OutputData: httpOutput, URL: urlsToCheck[i].URL, StatusCode: HttpStatusCode, Error: err.Error()})

						skyalerts.SendMail(d)
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

func RemoveOldEntry(db *gorm.DB, endpoints models.UserInput) {
	var urlCheck []models.AllEndpoints
	db.Find(&urlCheck)
	var urlsId models.AllEndpoints

	for i := 0; i < len(urlCheck); i++ {
		var findCount int
		for m := range endpoints.Domains {
			if endpoints.Domains[m].Name == urlCheck[i].URL {
				log.Println("found URL", urlCheck[i].URL)
				if urlCheck[i].Active != true {
					db.Model(&urlsId).Where("url = ?", urlCheck[i].URL).Update("Active", true)
				}
			} else {
				findCount = findCount + 1
			}
		}
		if findCount >= len(endpoints.Domains) {
			log.Println(urlCheck[i].URL, "Not found in input.json, removing the check")
			db.Model(&urlsId).Where("url = ?", urlCheck[i].URL).Update("Active", false)
		}
	}

}
