package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"somethingof/Api"
	"somethingof/Config"
	"somethingof/Models"
	"somethingof/Utils"
	"time"

	"google.golang.org/api/compute/v1"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func main() {
	// Set your GCP project, zone, and VM instance name
	projectID := "agentless-vasavi"
	labelKey := Config.UPTYCS_SCANNER_LABEL_KEY
	labelValue := Config.UPTYCS_SCANNER_LABEL_VALUE

	// Create a new context
	ctx := context.Background()

	// Create a new Compute Engine service using application credentials
	computeService, err := compute.NewService(ctx)
	if err != nil {
		log.Fatalf("Failed to create compute service: %v", err)
	}

	globalStateDict := make(map[string][]Models.DisksAtScanner)
	alsTracker := []Models.AlsTrackerObject{}
	fileCounter := 0

	for i := 0; i < 1; i++ {
		alsTracker = append(alsTracker, Models.AlsTrackerObject{})
		alsTracker[i].Time = time.Now()

		// Get instance IDs
		InstanceIDsList := Api.GetAllScannerIds(computeService, labelKey, labelValue, ctx, projectID)
		stateDict := make(map[string][]Models.DisksAtScanner)

		// Loop through instance IDs and get disks attached
		for _, instId := range InstanceIDsList {
			stateDict[instId.Name] = append(stateDict[instId.Name], Api.DisksAttachedToScanner(computeService, projectID, instId.Name, instId.Zone, ctx))
			globalStateDict[instId.Name] = append(globalStateDict[instId.Name], Api.DisksAttachedToScanner(computeService, projectID, instId.Name, instId.Zone, ctx))
		}
		alsTracker[i].DisksAtScanner = marshalDisksAtScanner(stateDict) // Marshal the map to JSON

		// Get ALS disks and snapshots
		alsTracker[i].AlsDisks = marshalDisksAtScanner(Api.DiskCreatedByAls(computeService, ctx, projectID))
		alsTracker[i].AlsSnapshots = marshalDisksAtScanner(Api.SnapshotsCreatedByALS(computeService, ctx, projectID))

		// Save to DB
		db := Config.GetDB()
		if err := storeAlsTracker(db, alsTracker); err != nil {
			log.Fatalf("Failed to create AlsTrackerDBObject: %v", err)
		}

		// Save JSON files
		if fileCounter%1 == 0 {
			filename := fmt.Sprintf("./DataFiles/Tracker/trackerObject-%d.json", fileCounter)
			Utils.CreateFile(filename, alsTracker)

			filename = fmt.Sprintf("./DataFiles/ScannerState/StateDictObject-%d.json", fileCounter)
			Utils.CreateFile(filename, globalStateDict)
		}

		fileCounter += 1
		fmt.Println("alsTracker: \n", alsTracker)
	}
}

// marshalDisksAtScanner converts the map to JSON for database storage
func marshalDisksAtScanner(stateDict interface{}) datatypes.JSON {
	jsonData, err := json.Marshal(stateDict)
	if err != nil {
		log.Fatalf("Failed to marshal DisksAtScanner: %v", err)
	}
	return datatypes.JSON(jsonData)
}

// storeAlsTracker saves the alsTracker object into the database
func storeAlsTracker(db *gorm.DB, alsTracker []Models.AlsTrackerObject) error {
	for i := range alsTracker {
		// Save the record to the DB
		if err := db.Table("als_table").Create(&alsTracker[i]).Error; err != nil {
			return err
		}
	}
	return nil
}
