package logger

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func New(prefix string, caller bool) *log.Logger {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		TimeFormat:      time.Kitchen,
		Prefix:          prefix,
		ReportCaller:    caller,
		ReportTimestamp: true,
	})

	return logger
}
