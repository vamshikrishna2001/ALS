package Models

import (
	"somethingof/Config"

	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = Config.GetDB()
	db.AutoMigrate(&AlsTrackerObject{})
}
