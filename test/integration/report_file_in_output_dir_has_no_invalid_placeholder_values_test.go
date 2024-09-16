package integration

import (
	"os"
	"strings"
	"testing"
)

func TestReportFileHasNoErrorPrVariableNotSetOrHasNoValue(t *testing.T) {
	reportFile := "VerificationReport_validation_Dummy1234567890_dummy_environment_name.html"
	outputDir := "output"
	filePath := "../../" + outputDir + "/" + reportFile

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file %q: %v", filePath, err)
	}

	if strings.Contains(string(content), "'pr' variable is not set or has no value") {
		t.Errorf("File %q contains invalid placeholder value \"'pr' variable is not set or has no value\"", filePath)
	}
}
