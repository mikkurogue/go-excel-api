package handlers

import (
	"encoding/json"
	"go-backend/core"
	"go-backend/util"
	"net/http"

	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func GetAllProcesses(c echo.Context) error {
	files, err := util.ReadDir()
	if err != nil {
		color.Red("Error: " + err.Error())
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error_code": "0x009",
			"message":   "Output directory does not exist, this means there are no processes complete.",
		})
	}

	var processes []string
	for _, file := range files {
		processes = append(processes, file.Name())
	}

	if len(processes) == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"error_code": "0x019",
			"message":    "No processes found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"processes": processes,
	})

}

func GetProcessById(c echo.Context) error {

	process_id := c.Param("id")
	if process_id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error_code": "0x001",
			"message":    "processId is required",
		})
	}

	shipment_file, err := util.ReadFile(process_id, "shipments.json")
	if err != nil {
		color.Red("Error: " + err.Error())
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error_code": "0x011",
			"message":    "Shipments json file not found for process",
		})
	}

	asset_file, err := util.ReadFile(process_id, "assets.json")
	if err != nil {
		color.Red("Error: " + err.Error())
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error_code": "0x012",
			"message":    "Assets json file not found for process",
		})
	}

	var shipments []core.Shipment
	var assets []core.EmissionAsset

	err = json.Unmarshal(shipment_file, &shipments)
	if err != nil {
		color.Red("Error unmarshalling shipments: " + err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error_ode": "0x013",
			"message":   "Error unmarshalling shipments json",
		})
	}

	err = json.Unmarshal(asset_file, &assets)
	if err != nil {
		color.Red("Error unmarshalling assets: " + err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error_ode": "0x014",
			"message":   "Error unmarshalling assets json",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"shipments": shipments,
		"assets":    assets,
	})

}
