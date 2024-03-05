package unit

import (
	"os/exec"
	"strings"
	"testing"
)

func TestGenerateVerificationReportArtifactName(t *testing.T) {
	tests := []struct {
		readyFor string
		expected string
	}{
		{"production", "verification_report_validation"},
		{"use", "verification_report_production"},
	}

	for _, test := range tests {
		t.Run(test.readyFor, func(t *testing.T) {
			cmd := exec.Command("bash", "-c", "../../src/script/get_verification_report_artifact_name_for_context.sh "+test.readyFor)
			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}

			actual := strings.TrimSpace(string(output))
			if actual != test.expected {
				t.Errorf("Expected %q, but got %q", test.expected, actual)
			}
		})
	}
}
