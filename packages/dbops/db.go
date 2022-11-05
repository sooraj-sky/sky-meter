package dbops

import (
	"encoding/json"
	"gorm.io/gorm"
	models "sky-meter/models"
	httpreponser "sky-meter/packages/httpres"
)

func InitialMigration(db *gorm.DB) {
	db.AutoMigrate(&models.AllEndpoints{})
	db.AutoMigrate(&models.HttpOutput{})
}

func InsertUrlsToDb(db *gorm.DB, endpoints models.JsonInput) {
	var urlCheck models.AllEndpoints
	for i := 0; i < len(endpoints); i++ {
		db.Where("URL=?", endpoints[i].URL).Find(&urlCheck)
		if urlCheck.CreatedAt == 0 && urlCheck.URL != endpoints[i].URL {
			db.Create(&models.AllEndpoints{URL: endpoints[i].URL, Timeout: endpoints[i].Timeout, SkipSsl: endpoints[i].SkipSsl, Frequency: endpoints[i].Frequency, Group: endpoints[i].Group})

		}

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
			httpOutput := httpreponser.CallEndpoint(urlsToCheck[i].URL)
			var byteHttpOutput models.Debug
			json.Unmarshal(httpOutput, &byteHttpOutput)
			db.Create(&models.HttpOutput{OutputData: httpOutput})
		} else {
			db.Model(&urlsId).Where("id = ?", urlsToCheck[i].ID).Update("next_run", urlsToCheck[i].NextRun-1)
		}
	}
}

func InputHttpToDb(db *gorm.DB) {

}
