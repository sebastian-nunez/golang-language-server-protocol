package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileLogger(t *testing.T) {
	// Arrange
	filename := "test.log"
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("unable to get current working directory: %v", err)
	}

	// Act
	logDir := filepath.Join(cwd, logsDir)
	logFilePath := filepath.Join(logDir, filename)
	os.RemoveAll(logDir) // Clean up logs directory before the test

	// Assert
	logger := NewFileLogger(filename)
	if logger == nil {
		t.Fatalf("expected logger to be non-nil")
	}

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		t.Fatalf("expected log file to be created at %s, but it does not exist", logFilePath)
	}

	fileInfo, err := os.Stat(logFilePath)
	if err != nil {
		t.Fatalf("unable to stat log file: %v", err)
	}

	expectedMode := os.FileMode(0666)
	actualMode := fileInfo.Mode().Perm() // Get only the permission bits
	if actualMode != expectedMode&actualMode {
		t.Errorf("expected log file to have mode %v, but got %v", expectedMode, actualMode)
	}

	os.RemoveAll(logDir)
}
