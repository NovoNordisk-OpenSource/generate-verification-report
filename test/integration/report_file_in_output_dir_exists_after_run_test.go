package integration

import (
	"os"
	"testing"
)

func TestReportFileExists(t *testing.T) {
	// Set up the test case
	reportFile := "VerificationReport_production_Dummy1234567890_dummy_environment_name.html"
	outputDir := "output"
	filePath := "../../" + outputDir + "/" + reportFile

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("File %q does not exist", filePath)
	}
}
