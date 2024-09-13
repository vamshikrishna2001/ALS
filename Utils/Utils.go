package Utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"somethingof/Config"
	"somethingof/Models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func CreateFile(Filename string, object interface{}) error {
	file, err := os.Create(Filename)
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(object, " ", "  ") // Indented for better readability
	if err != nil {
		fmt.Println("Error marshaling struct to JSON:", err)
	}

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
	return err

}

// storeAlsTracker saves the alsTracker object into the database
func StoreAlsTracker(db *gorm.DB, alsTracker Models.AlsTrackerObject) error {
	// for i := range alsTracker {
	// Save the record to the DB
	if err := db.Table(Config.TABLE_NAME).Create(&alsTracker).Error; err != nil {
		return err
		// }
	}
	return nil
}

// marshalDisksAtScanner converts the map to JSON for database storage
func MarshalDisksAtScanner(stateDict interface{}) datatypes.JSON {
	jsonData, err := json.Marshal(stateDict)
	if err != nil {
		log.Fatalf("Failed to marshal DisksAtScanner: %v", err)
	}
	return datatypes.JSON(jsonData)
}
