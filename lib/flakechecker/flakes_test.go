package flakechecker

import (
    "testing"
    "io/ioutil"
)

func TestLoadPatterns(t *testing.T) {
    patternsDir := "../../tests/testdata/testpatterns"
    expectedPatternCount := 4
    patterns, err := loadTestPatterns(t, patternsDir, expectedPatternCount)
    if err != nil {
        return
    }

    // Verify that SkipRetry defaults to false when not found in JSON
    for _, pattern := range patterns {
        if pattern.SkipRetry == true {
            t.Errorf("SkipRetry should default to false, but got %v", pattern.SkipRetry)
        }
    }
}

func TestCheckIfFlakeOccurred(t *testing.T) {

    patternsDir := "../../tests/testdata/buildlog"
    inputFile := "../../tests/testdata/buildlog/build-log.txt"
    expectedPatternCount := 3
    expectedFlakePatternCount := 2
    _, err := loadTestPatterns(t, patternsDir, expectedPatternCount)
    if err != nil {
        return
    }

    inputBytes, err := ioutil.ReadFile(inputFile)
    if err != nil {
        t.Fatalf("Error reading input file: %v", err)
    }
    input := string(inputBytes)

    flakePatterns, shouldRetry, err := CheckIfFlakeOccurred(input, patternsDir)
    if err != nil {
        t.Fatalf("Error occurred: %v", err)
    }

    if len(flakePatterns) != expectedFlakePatternCount {
        t.Errorf("Expected %d flake patterns, but got %d", expectedFlakePatternCount, len(flakePatterns))
    }

    expectedPatterns := []string{
		"Failed to check and update snapshot content: failed to remove VolumeSnapshotBeingCreated annotation on the content snapcontent-",
		"Backup and restore tests Backup and restore applications [It] MySQL application two Vol CSI",
	}
    for _, pattern := range flakePatterns {
        if !containsString(expectedPatterns, pattern.StringSearchPattern) {
            t.Errorf("Unexpected flake pattern found: %s", pattern.StringSearchPattern)
        }
    }
    if !shouldRetry {
        t.Errorf("Expected shouldRetry to be true, but got false")
    }
}

func TestCheckIfFlakeSkipOccurred(t *testing.T) {

    patternsDir := "../../tests/testdata/skipretry"
    inputFile := "../../tests/testdata/buildlog/build-log.txt"
    expectedPatternCount := 3
    expectedFlakePatternCount := 2
    _, err := loadTestPatterns(t, patternsDir, expectedPatternCount)
    if err != nil {
        return
    }

    inputBytes, err := ioutil.ReadFile(inputFile)
    if err != nil {
        t.Fatalf("Error reading input file: %v", err)
    }
    input := string(inputBytes)

    flakePatterns, shouldRetry, err := CheckIfFlakeOccurred(input, patternsDir)
    if err != nil {
        t.Fatalf("Error occurred: %v", err)
    }

    if len(flakePatterns) != expectedFlakePatternCount {
        t.Errorf("Expected %d flake patterns, but got %d", expectedFlakePatternCount, len(flakePatterns))
    }

    expectedPatterns := []string{
		"Failed to check and update snapshot content: failed to remove VolumeSnapshotBeingCreated annotation on the content snapcontent-",
		"Backup and restore tests Backup and restore applications [It] MySQL application two Vol CSI",
	}
    for _, pattern := range flakePatterns {
        if !containsString(expectedPatterns, pattern.StringSearchPattern) {
            t.Errorf("Unexpected flake pattern found: %s", pattern.StringSearchPattern)
        }
    }
    if !shouldRetry {
        t.Errorf("Expected shouldRetry to be true, but got false")
    }
}

func TestCheckIfFlakeOccurredLoadDefaultPatterns(t *testing.T) {
    flakePatterns, _, err := CheckIfFlakeOccurred("")
    if err != nil {
        t.Fatalf("Error occurred: %v", err)
    }

    expectedPatternCount := 0
    if len(flakePatterns) != expectedPatternCount {
        t.Errorf("Expected %d default flake patterns, but got %d", expectedPatternCount, len(flakePatterns))
    }
}

func loadTestPatterns(t *testing.T, patternsDir string, expectedPatternCount int) ([]FlakePattern, error) {
    patterns, err := loadPatterns(patternsDir)
    if err != nil {
        t.Fatalf("Error loading patterns: %v", err)
    }

    // Verify that patterns are loaded correctly
    if len(patterns) != expectedPatternCount {
        t.Errorf("Expected %d flake patterns, but got %d", expectedPatternCount, len(patterns))
    }

    return patterns, nil
}

func containsString(slice []string, str string) bool {
    for _, s := range slice {
        if s == str {
            return true
        }
    }
    return false
}