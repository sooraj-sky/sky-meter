package dbops

import (
	"encoding/json"
	"log"
	models "sky-meter/models"
	httpreponser "sky-meter/packages/httpres"

	"gorm.io/gorm"
)

type error interface {
	Error() string
}

func InitialMigration(db *gorm.DB) {
	db.AutoMigrate(&models.AllEndpoints{})
	db.AutoMigrate(&models.HttpOutput{})
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

	db.Find(&urlsToCheck)
	for i := 0; i < len(urlsToCheck); i++ {
		db.First(&urlsId, urlsToCheck[i].ID)
		if urlsToCheck[i].NextRun == 0 {
			db.Model(&urlsId).Where("id = ?", urlsToCheck[i].ID).Update("next_run", urlsToCheck[i].Frequency)
			httpOutput, HttpStatusCode, err := httpreponser.CallEndpoint(urlsToCheck[i].URL, urlsToCheck[i].Timeout, urlsToCheck[i].SkipSsl)
			if err != nil {
				db.Create(&models.HttpOutput{OutputData: httpOutput, URL: urlsToCheck[i].URL, StatusCode: HttpStatusCode, Error: err.Error()})
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
