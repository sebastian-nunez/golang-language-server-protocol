package util

import (
	"log"
	"os"
	"path/filepath"
)

const logsDir = "logs"

// GetLogger returns a logger that writes to the specified file. Logs are written to
// to the current working directory inside a `logs` folder.
func NewFileLogger(filename string) *log.Logger {
	cwd, err := os.Getwd()
	if err != nil {
		panic("unable to get current working directory: " + err.Error())
	}

	logDir := filepath.Join(cwd, logsDir)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("unable to create directory for log file: " + err.Error())
	}
	logFilePath := filepath.Join(logDir, filename)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("unable to open log file")
	}
	return log.New(file, "[golang-lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
