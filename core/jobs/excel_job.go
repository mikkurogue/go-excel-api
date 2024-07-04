package jobs

import (
	"encoding/json"
	"errors"
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

func (c *CoreJobExcel) Start(fileName string) (core.ParseError) {
	c.Started = true
	c.Finished = false
	data, err := core.ProcessExcel(fileName)

	if err.Error != nil {
		color.Red("Error processing file: %v", err.Error)
		log.Fatal(err.Error)
		// naked return here to abort the rest of the function
		return c.Abort("process aborted die to error in file processing")
	}

	c.Shipments = data.Shipments
	c.Assets = data.Assets

	c.ExportProcessJsonFiles()

	c.DeleteFileOnComplete(fileName)

	// test to return an empty error meaning its "fine" as the error should be nil
	return core.ParseError{}
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
	shipments, shipmentsError := json.MarshalIndent(c.Shipments, "", "    ")
	if shipmentsError != nil {
		color.Red("Error marshalling shipments to JSON: %v", shipmentsError)
	}

	// Marshal assets to JSON
	assets, assetsError := json.MarshalIndent(c.Assets, "", "    ")
	if assetsError != nil {
		color.Red("Error marshalling assets to JSON: %v", assetsError)
	}

	// Write shipments JSON to file
	shipmentsFilePath := filepath.Join(dir, "shipments.json")
	if shipmentsError := os.WriteFile(shipmentsFilePath, shipments, 0644); shipmentsError != nil {
		color.Red("Error writing shipments JSON to file '%s': %v", shipmentsFilePath, shipmentsError)
	}

	// Write assets JSON to file
	assetsFilePath := filepath.Join(dir, "assets.json")
	if assetsError := os.WriteFile(assetsFilePath, assets, 0644); assetsError != nil {
		color.Red("Error writing assets JSON to file '%s': %v", assetsFilePath, assetsError)
	}

	color.HiGreen("Process %s finished\n", c.ProcessId)
	c.Finished = true
	c.Started = false
}

/*
DeleteFileOnComplete deletes the file after the process is complete.
We do this to remove disk space usage and clutter - lets hope this wouldnt be "unit" costs on cloud providers.
*/
func (c *CoreJobExcel) DeleteFileOnComplete(fileName string) {
	e := os.Remove(fileName)
	if e != nil {
		color.Red("Error deleting file: %v", e)
		log.Fatal(e)
	}
}

func (c *CoreJobExcel) Abort(reason string) core.ParseError {
	return core.ParseError{
		Error: errors.New(reason),
	}
}