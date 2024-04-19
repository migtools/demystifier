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

package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

const logFile = "../../tests/testdata/buildlog/build-log.txt"

func TestPrintTestSummary(t *testing.T) {
	type args struct {
		logFile string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test with build-log.txt",
			args: args{
				logFile: logFile,
			},
			want: `Test Summary Table:
+-------------------------------------------------------------------------------+--------------+------------+------------------+
| TEST NAME                                                                     | NUM ATTEMPTS | NUM FAILED | AVERAGE RUN TIME |
+-------------------------------------------------------------------------------+--------------+------------+------------------+
| Should succeed                                                                |            1 |          0 |           5.044s |
| AWS Without Region And S3ForcePathStyle true should fail                      |            1 |          0 |          20.036s |
| Should succeed                                                                |            1 |          0 |          20.071s |
| HTTP_PROXY set                                                                |            1 |          0 |          35.243s |
| NO_PROXY set                                                                  |            1 |          0 |          35.291s |
| unsupportedOverrides should succeed                                           |            1 |          0 |        1m20.133s |
| Adding CSI plugin                                                             |            1 |          0 |        1m20.133s |
| Provider plugin                                                               |            1 |          0 |        1m20.136s |
| AWS With Region And S3ForcePathStyle should succeed                           |            1 |          0 |        1m20.138s |
| Adding Velero custom plugin                                                   |            1 |          0 |        1m20.139s |
| Set restic node selector                                                      |            1 |          0 |         1m20.14s |
| AWS Without Region No S3ForcePathStyle with BackupImages false should succeed |            1 |          0 |        1m20.141s |
| NoDefaultBackupLocation                                                       |            1 |          0 |        1m20.141s |
| Default velero CR, test carriage return                                       |            1 |          0 |        1m20.141s |
| Default velero CR                                                             |            1 |          0 |        1m20.142s |
| DPA CR with bsl and vsl                                                       |            1 |          0 |        1m20.143s |
| Enable tolerations                                                            |            1 |          0 |        1m20.148s |
| Adding Velero resource allocations                                            |            1 |          0 |        1m20.153s |
| Default velero CR with restic disabled                                        |            1 |          0 |        1m20.172s |
| HTTPS_PROXY set                                                               |            1 |          0 |         2m5.099s |
| Mongo application KOPIA                                                       |            1 |          0 |        2m31.823s |
| MySQL application KOPIA                                                       |            1 |          0 |         2m36.65s |
| MySQL application RESTIC                                                      |            1 |          0 |        2m46.649s |
| Mongo application RESTIC                                                      |            1 |          0 |        2m51.694s |
| MySQL application CSI                                                         |            2 |          1 |         3m8.943s |
| Config unset                                                                  |            1 |          0 |        3m31.199s |
| Mongo application CSI                                                         |            1 |          0 |        3m36.749s |
| MySQL application DATAMOVER                                                   |            1 |          0 |        4m16.949s |
| Mongo application DATAMOVER                                                   |            1 |          0 |        4m17.239s |
| Mongo application DATAMOVER                                                   |            1 |          0 |        4m36.933s |
| Mongo application BlockDevice DATAMOVER                                       |            1 |          0 |         5m6.999s |
| MySQL application two Vol CSI                                                 |            3 |          3 |        6m16.036s |
+-------------------------------------------------------------------------------+--------------+------------+------------------+
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testData, err := parseLogFile(tt.args.logFile)
			if err != nil {
				t.Errorf("Error parsing log file: %v", err)
			}
			old := os.Stdout // keep backup of the real stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PrintTestSummary(testData)

			outC := make(chan string)
			// copy the output in a separate goroutine so printing can't block indefinitely
			go func() {
				var buf bytes.Buffer
				io.Copy(&buf, r)
				outC <- buf.String()
			}()
			w.Close()
			os.Stdout = old
			outString := <-outC
			if outString != tt.want {
				t.Errorf("PrintTestSummary() = %v, want %v", outString, tt.want)
			}
		})
	}
}
