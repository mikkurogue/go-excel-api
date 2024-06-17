package handlers

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func GetAllProcesses(c echo.Context) error {
	files, err := ReadDir()
	if err != nil {
		color.Red("Error: " + err.Error())
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"errorCode": "0x009",
			"message":   "Output directory does not exist, this means there are no processes complete.",
		})
	}

	var processes []string
	for _, file := range files {
		processes = append(processes, file.Name())
	}

	if len(processes) == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "No processes found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"processes": processes,
	})

}

func ReadDir() ([]fs.DirEntry, error) {
	return os.ReadDir("./output")
}
