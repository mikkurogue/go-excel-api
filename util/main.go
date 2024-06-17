package util

import (
	"time"

	"github.com/fatih/color"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	color.Cyan("%s took %s", name, elapsed)
}
