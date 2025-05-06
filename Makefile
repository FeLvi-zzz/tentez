BINDIR:=dist
ROOT_PACKAGE:=$(shell go list .)
COMMAND_PACKAGES:=$(shell go list ./cmd/...)
BINARIES:=$(COMMAND_PACKAGES:$(ROOT_PACKAGE)/cmd/%=$(BINDIR)/%)

GOVERSION=$(shell go version)
CURRENT_GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
CURRENT_GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
GOOS=$(CURRENT_GOOS)
GOARCH=$(CURRENT_GOARCH)

BUILD_TARGETS= \
	build-linux-arm64 \
	build-linux-amd64 \
	build-darwin-amd64 \
	build-darwin-arm64 \
	build-windows-amd64

PLATFORMS=darwin linux windows
ARCHITECTURES=amd64 arm64

GO_FILES:=$(shell find . -type f -name '*.go' -print)

.PHONY: build build_all clean

lint: golangci-lint
	$(GOLANGCI_LINT) run

test: lint
	go test ./... -coverprofile cover.out

build_all: $(BUILD_TARGETS)

build-linux-arm64:
	@$(MAKE) build GOOS=linux GOARCH=arm64

build-linux-amd64:
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-darwin-arm64:
	@$(MAKE) build GOOS=darwin GOARCH=arm64

build-darwin-amd64:
	@$(MAKE) build GOOS=darwin GOARCH=amd64

build-windows-amd64:
	@$(MAKE) build GOOS=windows GOARCH=amd64 SUFFIX=.exe

build: $(BINARIES)

$(BINARIES): $(GO_FILES)
	CGO_ENABLED=0 go build -ldflags="-X github.com/FeLvi-zzz/tentez.Revision=$(shell git rev-parse --short HEAD)" -o $@-$(GOOS)-$(GOARCH)$(SUFFIX) $(@:$(BINDIR)/%=$(ROOT_PACKAGE)/cmd/%)

clean:
	rm $(BINARIES)*

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint

golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
