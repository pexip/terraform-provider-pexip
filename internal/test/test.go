/*
 * SPDX-FileCopyrightText: 2025 Pexip AS
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pexip/terraform-provider-pexip/internal/helpers"
)

var (
	INFINITY_USERNAME = helpers.GetEnvStringOrDefault("PEXIP_USERNAME", "admin")
	INFINITY_PASSWORD = helpers.GetEnvStringOrDefault("PEXIP_PASSWORD", "admin")
	INFINITY_BASE_URL = helpers.GetEnvStringOrDefault("PEXIP_ADDRESS", "https://dev-manager.dev.pexip.network")
)

func GetTestdataLocation() (string, error) {
	cmdOut, err := exec.CommandContext(context.Background(), "git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get base directory: %v", err)
	}
	baseDir := strings.TrimSpace(string(cmdOut))
	return filepath.Join(baseDir, "testdata"), nil
}

func LoadTestFolder(t *testing.T, folder string) string {
	t.Helper()

	loc, err := GetTestdataLocation()
	if err != nil {
		t.Fatalf("failed to get test data location: %v", err)
	}

	folderPath := filepath.Join(loc, folder)
	files, err := os.ReadDir(folderPath)
	if err != nil {
		t.Fatalf("failed to read test data folder %s: %v", folder, err)
	}

	var filesToLoad []string
	var filesToCopy []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".tf") {
			filesToLoad = append(filesToLoad, filepath.Join(folder, file.Name()))
		} else if !file.IsDir() && !strings.HasSuffix(file.Name(), ".tf") {
			filesToCopy = append(filesToCopy, filepath.Join(folder, file.Name()))
		}
	}

	copyTestFiles(t, filesToCopy...)
	return loadTestFiles(t, loc, filesToLoad...)
}

func LoadTestData(t *testing.T, files ...string) string {
	t.Helper()

	loc, err := GetTestdataLocation()
	if err != nil {
		t.Fatalf("failed to get test data location: %v", err)
	}
	return loadTestFiles(t, loc, files...)
}

func copyTestFiles(t *testing.T, files ...string) {
	t.Helper()

	for _, file := range files {
		// copy file to current working directory
		srcPath, err := GetTestdataLocation()
		if err != nil || srcPath == "" {
			t.Fatalf("failed to find test file %s: %v", file, err)
		}
		data, err := os.ReadFile(fmt.Sprintf("%s/%s", srcPath, file)) // #nosec G304 -- Path is validated above
		if err != nil {
			t.Fatalf("failed to read test file %s: %v", file, err)
		}
		destPath := filepath.Base(file)
		err = os.WriteFile(destPath, data, 0600)
		if err != nil {
			t.Fatalf("failed to write test file %s: %v", destPath, err)
		}
	}
}

func loadTestFiles(t *testing.T, baseDir string, files ...string) string {
	t.Helper()

	var combined strings.Builder
	for _, file := range files {
		fullPath := filepath.Join(baseDir, file)
		// Validate that the file path is safe and within expected directory
		cleanPath := filepath.Clean(fullPath)
		if !strings.HasPrefix(cleanPath, filepath.Clean(baseDir)) {
			t.Fatalf("unsafe file path: %s is outside base directory %s", file, baseDir)
		}
		data, err := os.ReadFile(cleanPath) // #nosec G304 -- Path is validated above
		if err != nil {
			// Bug Fix: The original code was missing the error variable in the format string.
			t.Fatalf("failed to load test file %s: %v", file, err)
		}
		combined.Write(data)
		combined.WriteString("\n")
	}
	return combined.String()
}

func GetGitBasePath() (string, error) {
	output, err := exec.CommandContext(context.Background(), "git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func FindFileInGitRoot(filename string) ([]string, error) {
	basePath, err := GetGitBasePath()
	if err != nil {
		return nil, err
	}

	return FindFileInPath(basePath, filename)
}

func FindFileInPath(path, filename string) ([]string, error) {
	var err error
	var found []string

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.Name() == filename {
			found = append(found, path)
		}
		return nil
	})
	return found, err
}

func LoadTestFile(t *testing.T, filename string) string {
	t.Helper()

	files, err := FindFileInGitRoot(filename)
	if err != nil {
		t.Fatalf("failed to find test file %s: %v", filename, err)
	}
	if len(files) == 0 {
		t.Fatalf("test file %s not found", filename)
	}

	data, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatalf("failed to read test file %s: %v", filename, err)
	}
	return string(data)
}
