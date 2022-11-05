package dbops

import (
	"gorm.io/gorm"
	"log"
	models "sky-meter/models"
	httpreponser "sky-meter/packages/httpres"
)

func InitialMigration(db *gorm.DB) {
	db.AutoMigrate(&models.AllEndpoints{})
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
	db.Find(&urlsToCheck)
	for i := 0; i < len(urlsToCheck); i++ {
		log.Println(urlsToCheck[i], "here", urlsToCheck[i].NextRun, urlsToCheck[i].ID, "freq", urlsToCheck[i].Frequency)
		if urlsToCheck[i].NextRun == 0 {
			db.Model(&urlsToCheck).Where("id = ?", urlsToCheck[i].ID).Where("url = ?", urlsToCheck[i].URL).Update("next_run", urlsToCheck[i].Frequency)
			httpreponser.CallEndpoint(urlsToCheck[i].URL)
		} else {
			db.Model(&urlsToCheck).Where("id = ?", urlsToCheck[i].ID).Where("url = ?", urlsToCheck[i].URL).Update("next_run", (urlsToCheck[i].NextRun)-1)
		}

	}
}
