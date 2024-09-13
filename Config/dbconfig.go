package Config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

// Initialize the database connection
func init() {
	dsn := "host=localhost user=postgres dbname=postgres password=mysecretpassword sslmode=disable"
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	db = d
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return db
}
