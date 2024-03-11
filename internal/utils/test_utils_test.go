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

package utils

import "testing"

func TestGenerateLogsURL(t *testing.T) {
	type args struct {
		originalURL string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test with pull-ci-openshift-oadp-operator-master-4.12-e2e-test-azure",
			args: args{
				originalURL: "https://prow.ci.openshift.org/view/gs/test-platform-results/pr-logs/pull/openshift_oadp-operator/1330/pull-ci-openshift-oadp-operator-master-4.12-e2e-test-azure/1757841602983759872",
			},
			want: "https://gcsweb-ci.apps.ci.l2s4.p1.openshiftapps.com/gcs/test-platform-results/pr-logs/pull/openshift_oadp-operator/1330/pull-ci-openshift-oadp-operator-master-4.12-e2e-test-azure/1757841602983759872/artifacts/e2e-test-azure/e2e/build-log.txt",
		},
		{
			name: "Test with pull-ci-openshift-oadp-operator-master-4.14-e2e-test-aws",
			args: args{
				originalURL: "https://prow.ci.openshift.org/view/gs/test-platform-results/pr-logs/pull/openshift_oadp-operator/1330/pull-ci-openshift-oadp-operator-master-4.14-e2e-test-aws/1757841603164114944",
			},
			want: "https://gcsweb-ci.apps.ci.l2s4.p1.openshiftapps.com/gcs/test-platform-results/pr-logs/pull/openshift_oadp-operator/1330/pull-ci-openshift-oadp-operator-master-4.14-e2e-test-aws/1757841603164114944/artifacts/e2e-test-aws/e2e/build-log.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GeneratesLogURL(tt.args.originalURL); got != tt.want {
				t.Errorf("GeneratesLogURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
