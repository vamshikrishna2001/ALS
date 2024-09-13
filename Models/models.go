package Models

import (
	"time"

	"github.com/jinzhu/gorm"
	"gorm.io/datatypes"
)

// DisksAtScanner represents the state of disks at a given time.
type DisksAtScanner struct {
	Time         time.Time     `json:"time"`
	ScannerState string        `json:"scannerstate"`
	Disks        []DiskDetails `json:"disks"`
}

// AlsSnapshots holds information about snapshots and their count.
type AlsSnapshots struct {
	Snapshots    []string `json:"snapshots"`
	NumSnapshots int      `json:"num_snapshots"`
}

// AlsDisks represents the details of disks and their count.
type AlsDisks struct {
	Disks    []string `json:"disks"`
	NumDisks int      `json:"num_disks"`
}

// AlsTrackerObject contains all relevant tracking information.
type AlsTrackerObject struct {
	gorm.Model
	Time time.Time `json:"time"`
	// DisksAtScanner map[string][]DisksAtScanner `json:"disks_at_scanner"`
	DisksAtScanner datatypes.JSON `json:"disks_at_scanner"`
	AlsSnapshots   datatypes.JSON `json:"als_snapshots"`
	AlsDisks       datatypes.JSON `json:"als_disks"`
}

// type AlsTrackerDBObject struct {
// 	gorm.Model
// 	Time time.Time `json:"time"`
// 	// DisksAtScanner []map[string][]DisksAtScanner `json:"disks_at_scanner"`
// 	DisksAtScanner datatypes.JSON `json:"disks_at_scanner"`
// 	AlsSnapshots   AlsSnapshots   `json:"als_snapshots"`
// 	AlsDisks       AlsDisks       `json:"als_disks"`
// }

// scanner object
type ScannerDetails struct {
	Name string
	Zone string
}

// disk object
type DiskDetails struct {
	Name string
	Size int64
}

func (AlsTrackerObject) TableName() string {
	return "als_table"
}
