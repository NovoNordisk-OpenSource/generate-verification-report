package unit

import (
	"os/exec"
	"testing"
)

func TestExtractRequirementsMapping(t *testing.T) {
	// Set up the test case
	cmd := exec.Command("python3", "../../src/script/extract_requirements_name_to_id_mapping.py", "-folder", "../integration/requirements")
	output, err := cmd.Output()
	if err != nil {
		t.Errorf("Failed to execute Python file: %v", err)
	}

	// Get the expected output
	expected := "{'Upper case a string': 'upper_case_feature', 'Reverse String': 'reverse_string_feat', 'day_of_week_feature': 'day_of_week_feature'}\n"

	// Compare the output
	actual := string(output)
	if actual != expected {
		t.Errorf("Expected %q, but got %q", expected, actual)
	}
}
