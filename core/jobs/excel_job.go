package jobs

import (
	"encoding/json"
	"go-backend/core"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

type CoreJobExcelExecutor interface {
	Start(fileName string)
	AssignProcessId()
	ExportProcessJsonFiles()
}

type CoreJobExcel struct {
	CoreJobExcelExecutor
	ProcessId string
	Started   bool
	Finished  bool
	Shipments []core.Shipment
	Assets    []core.EmissionAsset
}

func (c *CoreJobExcel) Start(fileName string) {
	c.Started = true
	c.Finished = false
	data := core.ProcessExcel(fileName)

	c.Shipments = data.Shipments
	c.Assets = data.Assets

	c.ExportProcessJsonFiles()
}

func (c *CoreJobExcel) AssignProcessId() {
	c.ProcessId = uuid.New().String()
}

func (c *CoreJobExcel) ExportProcessJsonFiles() {
	// Create the directory if it doesn't exist
	dir := filepath.Join("output", c.ProcessId)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		color.Red("Error creating directory '%s': %v", dir, err)
	}

	// Marshal shipments to JSON
	shipments, shipmentsErr := json.MarshalIndent(c.Shipments, "", "    ")
	if shipmentsErr != nil {
		color.Red("Error marshalling shipments to JSON: %v", shipmentsErr)
	}

	// Marshal assets to JSON
	assets, assetsErr := json.MarshalIndent(c.Assets, "", "    ")
	if assetsErr != nil {
		color.Red("Error marshalling assets to JSON: %v", assetsErr)
	}

	// Write shipments JSON to file
	shipmentsFilePath := filepath.Join(dir, "shipments.json")
	if shipmentsErr := os.WriteFile(shipmentsFilePath, shipments, 0644); shipmentsErr != nil {
		color.Red("Error writing shipments JSON to file '%s': %v", shipmentsFilePath, shipmentsErr)
	}

	// Write assets JSON to file
	assetsFilePath := filepath.Join(dir, "assets.json")
	if assetsErr := os.WriteFile(assetsFilePath, assets, 0644); assetsErr != nil {
		color.Red("Error writing shipments JSON to file '%s': %v", shipmentsFilePath, shipmentsErr)
	}

	color.HiGreen("Process %s finished\n", c.ProcessId)
	c.Finished = true
	c.Started = false
}
