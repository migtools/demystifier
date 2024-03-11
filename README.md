# Demystifier

A simple CLI to gather PROW logs, and parse them for further processing.


## Getting Started

### Prerequisites
- go version v1.19.0+
- shell

### Building the tool
```sh
$ make build
$ ./demystifier --help
```

** Cleanup the build**
```sh
$ make clean
```

### Running

#### Gather summary information about the PROW job run

```sh
$ ./demystifier "${URL}"

# Example, with URL from the GitHub PR comment
$ ./demystifier https://prow.ci.openshift.org/view/gs/test-platform-results/pr-logs/pull/openshift_oadp-operator/1266/pull-ci-openshift-oadp-operator-master-4.13-e2e-test-azure/1767186600720076800

# Example, with URL pointing directly to the log file
$ ./demystifier https://gcsweb-ci.apps.ci.l2s4.p1.openshiftapps.com/gcs/test-platform-results/pr-logs/pull/openshift_oadp-operator/1266/pull-ci-openshift-oadp-operator-master-4.13-e2e-test-azure/1767186600720076800/artifacts/e2e-test-azure/e2e/build-log.txt
```

#### Gather logs from the PROW job run and store them in a local folder

```sh
$ ./demystifier -f OUTPUT_LOGS_DIR "${URL}"

# Example, with URL from the GitHub PR comment, to dump the logs into /tmp/logs_dir folder
$ ./demystifier -f /tmp/logs_dir https://prow.ci.openshift.org/view/gs/test-platform-results/pr-logs/pull/openshift_oadp-operator/1266/pull-ci-openshift-oadp-operator-master-4.13-e2e-test-azure/1767186600720076800
```

### Tests

To run unit tests, run
```sh
$ make test
```

## License

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
