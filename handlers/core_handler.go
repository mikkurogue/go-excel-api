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
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}

	// close file src
	defer src.Close()

	// destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// copy file
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	// start the process
	jobs := jobs.CoreJobExcel{}

	go jobs.Start(file.Filename)
	jobs.AssignProcessId()

	return c.JSON(http.StatusOK, map[string]any{
		"message":    "success",
		"filename":   file.Filename,
		"size":       file.Size,
		"process_id": jobs.ProcessId,
	},
	)
}
