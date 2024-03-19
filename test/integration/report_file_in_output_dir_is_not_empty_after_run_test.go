package integration

import (
	"os"
	"testing"
)

func TestReportFileIsNotEmpty(t *testing.T) {
	// Set up the test case
	reportFile := "VerificationReport_production_Dummy1234567890_dummy_environment_name.html"
	outputDir := "output"
	filePath := "../../" + outputDir + "/" + reportFile

	// Check if the file is not empty
	fileInfo, _ := os.Stat(filePath)
	if fileInfo.Size() == 0 {
		t.Errorf("File %q is empty", filePath)
	}
}
