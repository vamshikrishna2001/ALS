package Models

import (
	"somethingof/Config"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	db = Config.GetDB()
	db.AutoMigrate(&AlsTrackerObject{})
}
