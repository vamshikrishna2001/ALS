package Api

import (
	"context"
	"fmt"
	"log"
	"somethingof/Config"
	"somethingof/Models"
	"strings"
	"time"

	"google.golang.org/api/compute/v1"
)

func DisksAttachedToScanner(computeService *compute.Service, projectID string, instanceId string, instanceZone string, ctx context.Context) Models.DisksAtScanner {
	instance, err := computeService.Instances.Get(projectID, instanceZone, instanceId).Context(ctx).Do()

	if err != nil {
		fmt.Println("dsffv ", err)
	}
	DiskObjs := []Models.DiskDetails{}
	for _, disk := range instance.Disks {
		fmt.Printf("  - Disk Name: %s, Device Name: %s, Type: %s ,size: %d\n", disk.Source, disk.DeviceName, disk.Type, disk.DiskSizeGb)
		DiskObjs = append(DiskObjs, Models.DiskDetails{Name: disk.DeviceName, Size: disk.DiskSizeGb})

	}
	return Models.DisksAtScanner{ScannerState: instance.Status, Disks: DiskObjs, Time: time.Now()}

}

func SnapshotsCreatedByALS(computeService *compute.Service, ctx context.Context, projectID string) Models.AlsSnapshots {
	filter := fmt.Sprintf("labels.%s=%s", Config.ALS_LABEL_KEY, Config.ALS_LABEL_VALUE)

	snapshotsListCall := computeService.Snapshots.List(projectID).Filter(filter)
	snapshotsListResp, err := snapshotsListCall.Context(ctx).Do()
	if err != nil {
		log.Fatalf("Failed to list snapshots: %v", err)
	}
	snapsList := Models.AlsSnapshots{}
	// Print out the information about the snapshots with the given label
	snapsCounter := 0
	if len(snapshotsListResp.Items) == 0 {
		fmt.Println("No snapshots found with the specified label.")
	} else {
		fmt.Println("Snapshots with the specified label:")
		for _, snapshot := range snapshotsListResp.Items {
			fmt.Printf("Snapshot Name: %s, Creation Time: %s, Size: %d GB\n", snapshot.Name, snapshot.CreationTimestamp, snapshot.DiskSizeGb)
			snapsList.Snapshots = append(snapsList.Snapshots, snapshot.Name)
			snapsCounter += 1
		}
	}
	snapsList.NumSnapshots = snapsCounter
	return snapsList

}

func DiskCreatedByAls(computeService *compute.Service, ctx context.Context, projectID string) Models.AlsDisks {
	filter := fmt.Sprintf("labels.%s=%s", Config.ALS_LABEL_KEY, Config.ALS_LABEL_VALUE)

	// List all disks across all zones
	disksListCall := computeService.Disks.AggregatedList(projectID).Filter(filter)
	disksListResp, err := disksListCall.Context(ctx).Do()
	if err != nil {
		log.Fatalf("Failed to list disks: %v", err)
	}
	diskslist := Models.AlsDisks{}
	// Print out the information about the disks with the given label
	diskCount := 0
	for _, disksScopedList := range disksListResp.Items {
		if disksScopedList.Disks != nil {
			for _, disk := range disksScopedList.Disks {
				fmt.Printf("Disk Name: %s, Zone: %s, Size: %d GB\n", disk.Name, disk.Zone, disk.SizeGb)
				diskslist.Disks = append(diskslist.Disks, disk.Name)
				diskCount += 1
			}
		}
	}
	diskslist.NumDisks = diskCount
	return diskslist
}

func GetAllScannerIds(computeService *compute.Service, labelKey string, labelValue string, ctx context.Context, projectID string) []*Models.ScannerDetails {

	scd := []*Models.ScannerDetails{}
	// Define the label filter in the format "labels.key=value"

	filter := fmt.Sprintf("labels.%s=%s", labelKey, labelValue)

	// Use the AggregatedList method to get instances across all zones
	aggregatedListCall := computeService.Instances.AggregatedList(projectID).Filter(filter)
	aggregatedListResp, err := aggregatedListCall.Context(ctx).Do()
	if err != nil {
		log.Fatalf("Failed to retrieve instances: %v", err)
	}

	// Print out the instance details that match the label across all zones
	if len(aggregatedListResp.Items) == 0 {
		fmt.Println("No instances found with the specified label.")
	} else {
		// fmt.Println("ergreb", len(aggregatedListResp.Items))
		fmt.Printf("Instances with label %s=%s:\n", labelKey, labelValue)
		for _, instancesScopedList := range aggregatedListResp.Items {
			if instancesScopedList.Instances != nil {
				// fmt.Printf("Zone: %s\n", zone)
				for _, instance := range instancesScopedList.Instances {
					// fmt.Printf("  - Instance Name: %s, Instance ID: %d\n", instance.Name, instance.Id)
					parts := strings.Split(instance.Zone, "/")
					scd = append(scd, &Models.ScannerDetails{Name: instance.Name, Zone: strings.TrimPrefix(parts[len(parts)-1], "zones/")})

				}
			}
		}
	}
	return scd

}
