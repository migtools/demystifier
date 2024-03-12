package flakechecker

import (
    "encoding/json"
    "io/ioutil"
    "path/filepath"
    "strings"
)

// Pattern represents a single Flake
type FlakePattern struct {
    Issue               string `json:"issue"`
    Description         string `json:"description"`
    StringSearchPattern string `json:"string_search_pattern"`
    SkipRetry           bool   `json:"skip_retry,omitempty"`
}

// CheckIfFlakeOccurred checks if any flake patterns occurred in the given input string.
// It searches for patterns in the input string and returns a list of FlakePattern structs
// representing the flake patterns found, along with a boolean indicating if retry should occur.
//
// Parameters:
//   - input:       The input string to search for flake patterns.
//   - patternsDir: (Optional) The directory path containing the JSON files with flake patterns.
//                  If not provided, it defaults to "patterns".
//
// Returns:
//   - flakePatterns: A slice of FlakePattern structs representing the flake patterns found in the input string.
//   - shouldRetry:   A boolean indicating whether a retry should occur.
//   - err:           An error if any occurred during the process of loading patterns or searching for patterns.
//                    If no error occurred, err is nil.
func CheckIfFlakeOccurred(input string, patternsDir ...string) (flakePatterns []FlakePattern, shouldRetry bool, err error) {

    dir := "patterns"
    shouldRetry = false

    if len(patternsDir) > 0 {
        dir = patternsDir[0]
    }

    patterns, err := loadPatterns(dir)
    if err != nil {
        return nil, false, err
    }

    for _, pattern := range patterns {
        if strings.Contains(input, pattern.StringSearchPattern) {
            flakePatterns = append(flakePatterns, pattern)
            if !pattern.SkipRetry {
                shouldRetry = true
            }
        }
    }

    return flakePatterns, shouldRetry, nil
}

// loadPatterns loads flake patterns from JSON files in the specified subfolder.
// It recursively searches for JSON files in the given subfolder and its subdirectories,
// extracts flake patterns from each file, and returns a slice of all patterns found.
//
// Parameters:
//   - subfolderPath: A string representing the path to the folder containing JSON files
//                    with flake patterns. Subfolders are also searched recursively.
//
// Returns:
//   - patterns: A slice of FlakePattern structs representing flake patterns extracted
//               from the JSON files found in the specified subfolder and its subdirectories.
//   - err:      An error if any occurred during the process of loading patterns.
//               If no error occurred, err is nil.
func loadPatterns(subfolderPath string) (patterns []FlakePattern, err error) {
    patternFiles, err := filepath.Glob(filepath.Join(subfolderPath, "*.json"))
    if err != nil {
        return nil, err
    }

    for _, file := range patternFiles {
        patternData, err := ioutil.ReadFile(file)
        if err != nil {
            return nil, err
        }

        var filePatterns []FlakePattern
        err = json.Unmarshal(patternData, &filePatterns)
        if err != nil {
            return nil, err
        }

        for i := range filePatterns {
            if !filePatterns[i].SkipRetry {
                filePatterns[i].SkipRetry = false
            }
        }

        patterns = append(patterns, filePatterns...)
    }

    return patterns, nil
}