package dbops

import (
	"encoding/json"
	"log"

	models "github.com/sooraj-sky/sky-meter/models"
	skyalerts "github.com/sooraj-sky/sky-meter/packages/alerts"
	httpreponser "github.com/sooraj-sky/sky-meter/packages/httpres"

	"gorm.io/gorm"
)

type error interface {
	Error() string
}

func InitialMigration(db *gorm.DB) {
	db.AutoMigrate(&models.AllEndpoints{})
	db.AutoMigrate(&models.HttpOutput{})
	db.AutoMigrate(&models.OpsgenieAlertData{})
}

func InsertUrlsToDb(db *gorm.DB, endpoints models.JsonInput) {
	var urlCheck models.AllEndpoints
	var urlsId models.AllEndpoints
	for i := 0; i < len(endpoints); i++ {
		db.Where("URL=?", endpoints[i].URL).Find(&urlCheck)
		if urlCheck.CreatedAt == 0 && urlCheck.URL != endpoints[i].URL {
			db.Create(&models.AllEndpoints{URL: endpoints[i].URL, Timeout: endpoints[i].Timeout, SkipSsl: endpoints[i].SkipSsl, Frequency: endpoints[i].Frequency, Group: endpoints[i].Group, Active: true})
		}
		urlCheck = urlsId
	}
}

func GetUrlFrequency(db *gorm.DB) {
	var urlsToCheck []models.AllEndpoints
	var urlsId []models.AllEndpoints
	var alertStatus models.OpsgenieAlertData

	db.Find(&urlsToCheck)
	for i := 0; i < len(urlsToCheck); i++ {
		db.First(&urlsId, urlsToCheck[i].ID)
		if urlsToCheck[i].Active {
			if urlsToCheck[i].NextRun == 0 {
				db.Model(&urlsId).Where("id = ?", urlsToCheck[i].ID).Update("next_run", urlsToCheck[i].Frequency)
				httpOutput, HttpStatusCode, err := httpreponser.CallEndpoint(urlsToCheck[i].URL, urlsToCheck[i].Timeout, urlsToCheck[i].SkipSsl)
				if err != nil {
					db.First(&alertStatus, "url = ?", urlsToCheck[i].URL)

					if alertStatus.URL == urlsToCheck[i].URL {
						AlertStatus := skyalerts.CheckAlertStatus(alertStatus.RequestId)
						if (AlertStatus == "closed") || (alertStatus.Error != err.Error()) {
							//	alertReqId := skyalerts.OpsgenieCreateAlert(urlsToCheck[i].URL, err, urlsToCheck[i].Group)
							//	db.Model(&alertStatus).Where("url = ?", urlsToCheck[i].URL).Update("request_id", alertReqId)
							d := models.SmtpErr{"skywalks.in", "webiste is down again", "22-2-132", "unable to connect", 22, "onlyworkofficial@gmail.com", "linux.sooraj@gmail.com", "smtp.gmail.com", 587}

							skyalerts.SendMail(d)
						}

					} else {
						//	alertReqId := skyalerts.OpsgenieCreateAlert(urlsToCheck[i].URL, err, urlsToCheck[i].Group)
						//	db.Create(&models.OpsgenieAlertData{URL: urlsToCheck[i].URL, RequestId: alertReqId, Error: err.Error(), Active: true})
						//	db.Create(&models.HttpOutput{OutputData: httpOutput, URL: urlsToCheck[i].URL, StatusCode: HttpStatusCode, Error: err.Error()})
						d := models.SmtpErr{"skywalks.in", "webiste is down again", "22-2-132", "unable to connect", 22, "onlyworkofficial@gmail.com", "linux.sooraj@gmail.com", "smtp.gmail.com", 587}

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

func RemoveOldEntry(db *gorm.DB, endpoints models.JsonInput) {
	var urlCheck []models.AllEndpoints
	db.Find(&urlCheck)
	var urlsId models.AllEndpoints

	for i := 0; i < len(urlCheck); i++ {
		var findCount int
		for m := range endpoints {
			if endpoints[m].URL == urlCheck[i].URL {
				log.Println("found URL", urlCheck[i].URL)
				if urlCheck[i].Active != true {
					db.Model(&urlsId).Where("url = ?", urlCheck[i].URL).Update("Active", true)
				}
			} else {
				findCount = findCount + 1
			}
		}
		if findCount >= len(endpoints) {
			log.Println(urlCheck[i].URL, "Not found in input.json, removing the check")
			db.Model(&urlsId).Where("url = ?", urlCheck[i].URL).Update("Active", false)
		}
	}

}
