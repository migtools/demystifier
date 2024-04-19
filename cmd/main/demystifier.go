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

// Demystifier app
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/migtools/demystifier/lib/utils"
	log "github.com/sirupsen/logrus"
)

const saveFolderPerm = 0750

func parseLogFile(logFile string) (*utils.TestRunData, error) {
	testRunDataPtr, err := utils.GetRunDataFromLog(logFile)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error")
		return nil, err
	}

	err = utils.SetIndividualTestsFromLog(testRunDataPtr, "It")

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error")
		return nil, err
	}

	return testRunDataPtr, nil
}

// DumpTestsToFolder saves logs to a destination folder
func DumpTestsToFolder(testData *utils.TestRunData, folder string) {
	mkdirErr := os.MkdirAll(folder, saveFolderPerm)
	if mkdirErr != nil {
		log.WithFields(log.Fields{
			"error": mkdirErr,
		}).Fatal("Error creating dir")
	}
	for i := range testData.TestRun {
		thisRun := &testData.TestRun[i]
		for j := range thisRun.Attempt {
			thisAttempt := &thisRun.Attempt[j]
			err := thisAttempt.DumpLogsToFileWithPrefixes(j, folder, thisAttempt.Name, ": ")
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Fatal("Error dumping logs to file")
			}
		}
	}
}

// PrintTestSummary prints the summary of tests
func PrintTestSummary(testData *utils.TestRunData) {
	// Define a struct to hold the summary data
	type TestSummary struct {
		Name           string
		NumAttempts    int
		NumFailed      int
		TotalRunTime   time.Duration
		NumOver1Second int
		AverageRunTime time.Duration
	}

	// Initialize a slice to hold the summary data for each test run
	var summaries []TestSummary

	// Loop through each test run to collect summary data
	for i := range testData.TestRun {
		var numAttempts, failedAttempts, numOver1Second int
		totalRunTime := time.Duration(0)
		thisTest := &testData.TestRun[i]
		for j := range thisTest.Attempt {
			// Increment the number of attempts
			numAttempts++
			thisAttempt := &thisTest.Attempt[j]
			// If the attempt failed, increment the failed attempts counter
			if thisAttempt.Status.Status == "FAILED" {
				failedAttempts++
			}

			// If the duration is greater than 1 second, increment the counter
			if thisAttempt.Duration > time.Second {
				numOver1Second++
			}

			// Add the duration to the total run time
			totalRunTime += thisAttempt.Duration
		}

		// Calculate the average run time based on durations over 1 second
		var averageRunTime time.Duration
		if numOver1Second > 0 {
			averageRunTime = totalRunTime / time.Duration(numOver1Second)
		}

		// Append the summary data to the slice
		summaries = append(summaries, TestSummary{
			Name:           thisTest.ShortName,
			NumAttempts:    numAttempts,
			NumFailed:      failedAttempts,
			TotalRunTime:   totalRunTime,
			NumOver1Second: numOver1Second,
			AverageRunTime: averageRunTime,
		})
	}

	// Sort by avg time
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].AverageRunTime < summaries[j].AverageRunTime
	})

	// Print the summary table
	fmt.Println("Test Summary Table:")
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Test Name", "Num Attempts", "Num Failed", "Average Run Time"})
	for _, summary := range summaries {
		t.AppendRows([]table.Row{
			{summary.Name, summary.NumAttempts, summary.NumFailed, summary.AverageRunTime},
		})
	}
	t.Render()
}

func main() {
	log.SetLevel(log.InfoLevel)

	log.WithFields(log.Fields{
		">>> start_demystifier_timestamp": time.Now().Unix(),
	}).Info("Test Demystifier starts its journey")

	var (
		logLocation      string
		showPassing      bool
		timeStamps       bool
		debugMode        bool
		dumpLogsToFolder string
	)

	flag.BoolVar(&timeStamps, "t", false, "whether to include timestamps in the output (shorthand)")
	flag.BoolVar(&showPassing, "s", false, "show all tests even those passing")
	flag.BoolVar(&debugMode, "d", false, "debug mode")
	flag.StringVar(&dumpLogsToFolder, "f", "", "dump logs to folder")

	flag.Parse()

	if debugMode {
		log.SetLevel(log.DebugLevel)
	}

	if len(flag.Args()) > 0 {
		logLocation = utils.GeneratesLogURL(flag.Arg(0))
	}

	log.WithFields(log.Fields{
		">>> location": logLocation,
	}).Info("Using log from")

	testData, _ := parseLogFile(logLocation)

	for i := range testData.TestRun {
		failedAttempts := 0 // Initialize counter for failed attempts in this test run
		thisTest := &testData.TestRun[i]
		for j := range thisTest.Attempt {
			thisAttempt := &thisTest.Attempt[j]
			fields := log.Fields{
				"Name": thisTest.ShortName,
				"No":   thisAttempt.AttemptNo,
				"Time": thisAttempt.Duration,
			}

			// If the attempt failed or showPassing is true, log the attempt
			if thisAttempt.Status.Status == utils.Failed {
				log.WithFields(fields).Error("Failed attempt run")
				// Increment the counter if the attempt failed
				if thisAttempt.Status.Status == utils.Failed {
					failedAttempts++
				}
			} else if showPassing {
				log.WithFields(fields).Info("Pass attempt run")
			}
		}

		// Summary for this test run
		if failedAttempts > 0 {
			log.WithFields(log.Fields{
				"Name":   thisTest.Name,
				"Failed": failedAttempts,
			}).Info("Test Summary")
		}
	}
	if dumpLogsToFolder != "" {
		DumpTestsToFolder(testData, dumpLogsToFolder)
		os.Exit(0)
	}
	PrintTestSummary(testData)

	log.WithFields(log.Fields{
		">>> end_demystifier_timestamp": time.Now().Unix(),
	}).Info("Test Demystifier finishes its journey")
}
