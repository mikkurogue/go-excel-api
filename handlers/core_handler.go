package handlers

import (
	"go-backend/core/jobs"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// route is /core
func Core(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func UploadExcel(c echo.Context) error {

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "0x901",
			"message":    "no file found in request",
		},)
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "0x1",
			"message":    "server could not open file stream",
		},)
	}

	// close file src
	defer src.Close()

	// destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "0x2",
			"message":    "server could not create output directory for process",
		},)

	}
	defer dst.Close()

	// copy file
	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"error": "0x3",
			"message":    "server could not copy file to output directory",
		},)
	}

	// start the process
	jobs := jobs.CoreJobExcel{}

	// figure out why this needs to be in a routine
	go jobs.Start(file.Filename)
	
	jobs.AssignProcessId()

	return c.JSON(http.StatusOK, map[string]any{
		"message":    "success",
		"file_name":  file.Filename,
		"size":       file.Size,
		"process_id": jobs.ProcessId,
	},
	)
}

func DeleteProcess(c echo.Context) error {
	processId := c.Param("id")
	jobs := jobs.CoreJobExcel{}
	jobs.DeleteProcess(processId)
	return c.JSON(http.StatusOK, map[string]any{
		"message":    "success",
		"process_id": processId,
	},
	)
}