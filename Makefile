#
# Copyright 2024.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
## Location to install dependencies to

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GOBUILD = $(GOBIN)/go build
GOCLEAN = $(GOBIN)/go clean
GORUN = $(GOBIN)/go run
GOTEST = $(GOBIN)/go test

BINARY_NAME = demystifier

GOLANGCI_LINT = $(LOCALBIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
GOLANGCI_LINT_VERSION ?= v1.56.2

EC ?= $(LOCALBIN)/ec-$(EC_VERSION)
EC_VERSION ?= 2.8.0

.PHONY: editorconfig
editorconfig: $(LOCALBIN) ## Download editorconfig locally if necessary.
	@[ -f $(EC) ] || { \
	set -e ;\
	ec_binary=ec-$(shell go env GOOS)-$(shell go env GOARCH) ;\
	ec_tar=$(LOCALBIN)/$${ec_binary}.tar.gz ;\
	curl -sSLo $${ec_tar} https://github.com/editorconfig-checker/editorconfig-checker/releases/download/$(EC_VERSION)/$${ec_binary}.tar.gz ;\
	tar xzf $${ec_tar} ;\
	rm -rf $${ec_tar} ;\
	mv $(LOCALBIN)/$${ec_binary} $(EC) ;\
	}

.PHONY: all
all: build

# Build target
build:
	GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) ./cmd/main/demystifier.go
	chmod +x $(BINARY_NAME)

# Example make run ARGS="--help"
.PHONY: run
run:
	$(GORUN) ./cmd/main/demystifier.go $(ARGS)

# Clean target
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

.PHONY: test
test:
	$(GOTEST) $$(go list ./...)

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,${GOLANGCI_LINT_VERSION})

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter & yamllint
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary (ideally with version)
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f $(1) ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv "$$(echo "$(1)" | sed "s/-$(3)$$//")" $(1) ;\
}
endef

.PHONY: ec
ec: editorconfig ## Run file formatter checks against all project's files.
	$(EC)
