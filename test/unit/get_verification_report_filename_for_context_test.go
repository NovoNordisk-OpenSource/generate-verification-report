package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestGenerateVerificationReportFilename(t *testing.T) {
	// Set up the test case
	envName := "ramone-service1-val-eu-central1"
	buildID := "617829"
	readyFor := "use"

	// Run the bash script
	cmd := exec.Command("bash", "-c", "../../src/script/get_verification_report_filename_for_context.sh \""+envName+"\" \""+buildID+"\" \""+readyFor+"\"")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}

	// Get the expected output
	expected := "VerificationReport_production_617829_ramone_service1_val_eu_central1"

	// Compare the output
	actual := strings.TrimSpace(string(output))
	if actual != expected {
		t.Errorf("Expected %q, but got %q", expected, actual)
	}
}
