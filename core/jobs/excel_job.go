package jobs

import (
	"encoding/json"
	"go-backend/core"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

// Assign this structure to have specific methods and properties
type CoreJobExcel struct {
	CoreJobExcecutor
	JobStatus
	Shipments []core.Shipment
	Assets    []core.EmissionAsset
}

func (c *CoreJobExcel) Start(file_name string) {
	c.Started = true
	c.Finished = false
	data := core.ProcessExcel(file_name)

	c.Shipments = data.Shipments
	c.Assets = data.Assets

	c.ExportProcessJsonFiles()

	c.DeleteFileOnComplete(file_name)
}

func (c *CoreJobExcel) AssignProcessId() {
	c.ProcessId = uuid.New().String()
}

/*
ExportProcessJsonFiles exports the shipments and assets to JSON files
*/
func (c *CoreJobExcel) ExportProcessJsonFiles() {
	// Create the directory if it doesn't exist
	dir := filepath.Join("output", c.ProcessId)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		color.Red("Error creating directory '%s': %v", dir, err)
	}

	// Marshal shipments to JSON
	shipments, shipments_err := json.MarshalIndent(c.Shipments, "", "    ")
	if shipments_err != nil {
		color.Red("Error marshalling shipments to JSON: %v", shipments_err)
	}

	// Marshal assets to JSON
	assets, assets_err := json.MarshalIndent(c.Assets, "", "    ")
	if assets_err != nil {
		color.Red("Error marshalling assets to JSON: %v", assets_err)
	}

	// Write shipments JSON to file
	shipments_file_path := filepath.Join(dir, "shipments.json")
	if shipments_err := os.WriteFile(shipments_file_path, shipments, 0644); shipments_err != nil {
		color.Red("Error writing shipments JSON to file '%s': %v", shipments_file_path, shipments_err)
	}

	// Write assets JSON to file
	assets_file_path := filepath.Join(dir, "assets.json")
	if assets_err := os.WriteFile(assets_file_path, assets, 0644); assets_err != nil {
		color.Red("Error writing assets JSON to file '%s': %v", assets_file_path, assets_err)
	}

	color.HiGreen("Process %s finished\n", c.ProcessId)
	c.Finished = true
	c.Started = false
}

/*
DeleteFileOnComplete deletes the file after the process is complete.
We do this to remove disk space usage and clutter - lets hope this wouldnt be "unit" costs on cloud providers.
*/
func (c *CoreJobExcel) DeleteFileOnComplete(file_name string) {
	e := os.Remove(file_name)
	if e != nil {
		color.Red("Error deleting file: %v", e)
		log.Fatal(e)
	}
}
