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

	// channels for synchronising
	scannerState := make(chan map[string][]Models.DisksAtScanner)
	snapshotState := make(chan Models.AlsSnapshots)
	diskState := make(chan Models.AlsDisks)

	for i := 0; i < 1; i++ {
		// you have to initialise the slice with emypty values if not .. you cannot index the slice ... so we add the the empty object them over write it
		alsTracker = append(alsTracker, Models.AlsTrackerObject{})

		// Task ONE
		// Updating the time object
		alsTracker[i].Time = time.Now()

		// TASK TWO
		// Get instance IDs

		go func(computeService *compute.Service, ctx context.Context, projectID string, ch chan<- map[string][]Models.DisksAtScanner) {
			InstanceIDsList := Api.GetAllScannerIds(computeService, ctx, projectID)
			stateDict := make(map[string][]Models.DisksAtScanner)

			// Loop through instance IDs and get disks attached
			for _, instId := range InstanceIDsList {
				das := Api.DisksAttachedToScanner(computeService, projectID, instId.Name, instId.Zone, ctx) // das means disks attached to scanner
				stateDict[instId.Name] = append(stateDict[instId.Name], das)
				globalStateDict[instId.Name] = append(globalStateDict[instId.Name], das)
			}
			ch <- stateDict

		}(computeService, ctx, projectID, scannerState)

		// Task THREE AND FOUR
		// Get ALS disks and snapshots
		go func(computeService *compute.Service, ctx context.Context, projectID string, ch chan<- Models.AlsDisks) {
			ch <- Api.DiskCreatedByAls(computeService, ctx, projectID)
		}(computeService, ctx, projectID, diskState)

		go func(computeService *compute.Service, ctx context.Context, projectID string, ch chan<- Models.AlsSnapshots) {
			ch <- Api.SnapshotsCreatedByALS(computeService, ctx, projectID)
		}(computeService, ctx, projectID, snapshotState)

		alsTracker[i].AlsDisks = Utils.MarshalDisksAtScanner(<-diskState)
		alsTracker[i].AlsSnapshots = Utils.MarshalDisksAtScanner(<-snapshotState)
		alsTracker[i].DisksAtScanner = Utils.MarshalDisksAtScanner(<-scannerState) // Marshal the map to JSON

		// Save to DB
		db := Config.GetDB()
		if err := Utils.StoreAlsTracker(db, alsTracker[i]); err != nil {
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
