package dbops

import (
  "log"
  "gorm.io/gorm"
  models "sky-meter/models"
)


func InitialMigration(db *gorm.DB) {
	db.AutoMigrate(&models.AllEndpoints{})
}


func InsertUrlsToDb(db *gorm.DB,endpoints models.JsonInput) {
  var urlCheck models.AllEndpoints
  for i := 0; i < len(endpoints); i++ {
    db.Where("URL=?", endpoints[i].URL).Find(&urlCheck)
    log.Println(urlCheck.ID, urlCheck.CreatedAt)
    if urlCheck.CreatedAt == 0 {
      db.Create(&models.AllEndpoints{URL: endpoints[i].URL, Timeout: endpoints[i].Timeout, SkipSsl: endpoints[i].SkipSsl, Frequency: endpoints[i].Frequency, Group: endpoints[i].Group})

    }

  }

}
