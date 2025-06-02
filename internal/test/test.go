package test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

//nolint:unused
func GetTestdataLocation() (string, error) {
	cmdOut, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get base directory: %v", err)
	}
	return fmt.Sprintf("%s/testdata", strings.TrimSpace(string(cmdOut))), nil
}

//nolint:deadcode,unused
func LoadTestConfig(t *testing.T, file string) string {
	loc, err := GetTestdataLocation()
	if err != nil {
		t.Fatalf("failed to get test data location")
	}

	data, err := os.ReadFile(fmt.Sprintf("%s/%s", loc, file))
	if err != nil {
		t.Fatalf("failed to load test config %s", file)
	}
	return string(data)
}
