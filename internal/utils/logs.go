/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package utils for the demystifier CLI
package utils

import (
	"fmt"
	"os"
	"strings"
)

// DumpLogsToFileWithPrefixes saves logs as files
func (a *AttemptData) DumpLogsToFileWithPrefixes(attemptNo int, folder string, prefixes ...string) error {
	// replace / in name
	fileName := strings.ReplaceAll(a.Name, "/", "_")
	filename := fmt.Sprintf("%s/%s_%d.log", folder, fileName, attemptNo)

	file, err := os.Create(filename)

	if err != nil {
		return fmt.Errorf("error creating log file: %v", err)
	}
	defer file.Close()
	for i := range a.Logs {
		for j := range prefixes {
			if _, err := file.WriteString(prefixes[j]); err != nil {
				return err
			}
		}
		if _, err := file.WriteString(a.Logs[i] + "\n"); err != nil {
			return err
		}
	}
	return nil
}
