APP_VERSION ?= $(shell git describe --abbrev=5 --dirty --tags --always)
GIT_COMMIT := $(shell git rev-parse --short=8 HEAD)
BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

BINDIR := $(PWD)/bin
OUTPUT_DIR := $(PWD)/_output

GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH ?= amd64

LDFLAGS := $(LDFLAGS) -X github.com/bougou/alertmanager-webhook-adapter/pkg/version.Version=$(APP_VERSION)
LDFLAGS := $(LDFLAGS) -X github.com/bougou/alertmanager-webhook-adapter/pkg/version.Commit=$(GIT_COMMIT)
LDFLAGS := $(LDFLAGS) -X github.com/bougou/alertmanager-webhook-adapter/pkg/version.BuildAt=$(BUILD_TIME)

PATH := $(BINDIR):$(PATH)
SHELL := env PATH='$(PATH)' /bin/sh

all: build

# Run tests
test: fmt vet
	@# Disable --race until https://github.com/kubernetes-sigs/controller-runtime/issues/1171 is fixed.
	ginkgo --randomizeAllSpecs --randomizeSuites --failOnPending --flakeAttempts=2 \
			--cover --coverprofile cover.out --trace --progress  $(TEST_ARGS)\
			./pkg/... ./cmd/...

# Build alertmanager-webhook-adapter binary
build: fmt vet
	go build -ldflags "$(LDFLAGS)" -o $(OUTPUT_DIR)/alertmanager-webhook-adapter ./cmd/alertmanager-webhook-adapter

# Cross compiler
build-all: fmt vet
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -a -o $(OUTPUT_DIR)/alertmanager-webhook-adapter-$(APP_VERSION)-linux-amd64 ./cmd/alertmanager-webhook-adapter
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -a -o $(OUTPUT_DIR)/alertmanager-webhook-adapter-$(APP_VERSION)-linux-arm64 ./cmd/alertmanager-webhook-adapter
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -a -o $(OUTPUT_DIR)/alertmanager-webhook-adapter-$(APP_VERSION)-darwin-amd64 ./cmd/alertmanager-webhook-adapter
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -a -o $(OUTPUT_DIR)/alertmanager-webhook-adapter-$(APP_VERSION)-darwin-arm64 ./cmd/alertmanager-webhook-adapter

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

lint:
	$(BINDIR)/golangci-lint run --timeout 2m0s ./pkg/... ./cmd/...

dependencies:
	test -d $(BINDIR) || mkdir $(BINDIR)
	GOBIN=$(BINDIR) go install github.com/onsi/ginkgo/ginkgo@v1.16.4

	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $(BINDIR) latest
