package main

import (
	"context"
	"fmt"
	"log"
	"somethingof/Api"
	"somethingof/Config"
	"somethingof/Models"
	"somethingof/Utils"
	"time"

	"google.golang.org/api/compute/v1"
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

	for i := 0; i < 3; i++ {
		alsTracker = append(alsTracker, Models.AlsTrackerObject{})

		alsTracker[i].Time = time.Now()
		InstanceIDsList := Api.GetAllScannerIds(computeService, labelKey, labelValue, ctx, projectID)
		stateDict := make(map[string][]Models.DisksAtScanner)
		for _, instId := range InstanceIDsList {
			stateDict[instId.Name] = append(stateDict[instId.Name], Api.DisksAttachedToScanner(computeService, projectID, instId.Name, instId.Zone, ctx))
			globalStateDict[instId.Name] = append(globalStateDict[instId.Name], Api.DisksAttachedToScanner(computeService, projectID, instId.Name, instId.Zone, ctx))
			alsTracker[i].DisksAtScanner = append(alsTracker[i].DisksAtScanner, stateDict)

		}
		alsTracker[i].AlsDisks = Api.DiskCreatedByAls(computeService, ctx, projectID)
		alsTracker[i].AlsSnapshots = Api.SnapshotsCreatedByALS(computeService, ctx, projectID)

		if fileCounter%1 == 0 {
			filename := fmt.Sprintf("./DataFiles/Tracker/trackerObject-%d.json", fileCounter)
			Utils.CreateFile(filename, alsTracker)

			db := Config.GetDB()
			db.Create(alsTracker)

			filename = fmt.Sprintf("./DataFiles/ScannerState/StateDictObject-%d.json", fileCounter)
			Utils.CreateFile(filename, globalStateDict)

		}

		// time.Sleep(5 * time.Minute)
		fileCounter += 1

	}
	fmt.Println("Global ", globalStateDict)
	fmt.Println("state ", alsTracker)

	// Create a JSON file

}
