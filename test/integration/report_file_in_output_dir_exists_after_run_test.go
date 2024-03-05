package integration

import (
	"os"
	"testing"
)

func TestReportFileExists(t *testing.T) {
	// Set up the test case
	reportFile := "report.html"
	outputDir := "output"
	filePath := "../../" + outputDir + "/" + reportFile

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("File %q does not exist", filePath)
	}
}
