package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"somethingof/Api"
	"somethingof/Config"
	"somethingof/Models"
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

	}
	fmt.Println("Global ", globalStateDict)
	fmt.Println("state ", alsTracker)
	jsonData, err := json.MarshalIndent(alsTracker, " ", "  ") // Indented for better readability
	if err != nil {
		fmt.Println("Error marshaling struct to JSON:", err)
		return
	}

	// Create a JSON file
	file, err := os.Create("disk_info.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
