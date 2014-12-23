package utils

import (
	"os"

	"github.com/motemen/go-colorine"
)

var logger = &colorine.Logger{
	colorine.Prefixes{
		"verbose": colorine.Verbose,
		"notice":  colorine.Notice,
		"info":    colorine.Info,
		"warn":    colorine.Warn,
		"error":   colorine.Error,
	},
}

// Log outputs `message` with `prefix` by go-colorine
func Log(prefix, message string) {
	logger.Log(prefix, message)
}

// ErrorIf outputs log if `err` occurs.
func ErrorIf(err error) bool {
	if err != nil {
		Log("error", err.Error())
		return true
	}

	return false
}

// DieIf outputs log and exit(1) if `err` occurs.
func DieIf(err error) {
	if err != nil {
		Log("error", err.Error())
		os.Exit(1)
	}
}
