package main

import (
	"os"
	"testing"
)

func TestReportFileExists(t *testing.T) {
	// Set up the test case
	reportFile := "report.html"
	outputDir := "output"
	filePath := "../../src/" + outputDir + "/" + reportFile

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("File %q does not exist", filePath)
	}
}

func TestReportFileIsNotEmpty(t *testing.T) {
	// Set up the test case
	reportFile := "report.html"
	outputDir := "output"
	filePath := "../../src/" + outputDir + "/" + reportFile

	// Check if the file is not empty
	fileInfo, _ := os.Stat(filePath)
	if fileInfo.Size() == 0 {
		t.Errorf("File %q is empty", filePath)
	}
}
