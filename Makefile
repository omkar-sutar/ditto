# Makefile for Ditto file transfer application
# Supports building for both Windows and Linux

# Binary names
BINARY_NAME_LINUX=ditto
BINARY_NAME_WINDOWS=ditto.exe

# Version information
VERSION=1.0.0
BUILD=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD) -X main.BuildDate=$(BUILD_DATE)"

# Different platforms
PLATFORMS=linux windows

# Source files
SRC_FILES=main.go client.go server.go control.go

# Default target
.PHONY: all
all: build

# Build for current platform
.PHONY: build
build:
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME_LINUX) $(SRC_FILES)
	@echo "Build successful for current platform"

# Clean build files
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf bin/
	@echo "Cleaned build files"

# Create build directory
.PHONY: prepare
prepare:
	mkdir -p bin

# Check dependencies
.PHONY: deps
deps:
	$(GOMOD) tidy
	@echo "Dependencies updated"

# Test the code
.PHONY: test
test:
	$(GOTEST) -v ./...

# Cross-compilation for all platforms
.PHONY: cross-build
cross-build: prepare linux windows
	@echo "Cross-compilation complete"

# Build for Linux
.PHONY: linux
linux: prepare
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/linux_amd64/$(BINARY_NAME_LINUX) $(SRC_FILES)
	GOOS=linux GOARCH=386 $(GOBUILD) $(LDFLAGS) -o bin/linux_386/$(BINARY_NAME_LINUX) $(SRC_FILES)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o bin/linux_arm64/$(BINARY_NAME_LINUX) $(SRC_FILES)
	@echo "Linux builds completed"

# Build for Windows
.PHONY: windows
windows: prepare
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/windows_amd64/$(BINARY_NAME_WINDOWS) $(SRC_FILES)
	GOOS=windows GOARCH=386 $(GOBUILD) $(LDFLAGS) -o bin/windows_386/$(BINARY_NAME_WINDOWS) $(SRC_FILES)
	@echo "Windows builds completed"

# Create distribution packages
.PHONY: dist
dist: cross-build
	mkdir -p dist
	# Linux packages
	cd bin/linux_amd64 && tar -czf ../../dist/$(BINARY_NAME_LINUX)_$(VERSION)_linux_amd64.tar.gz $(BINARY_NAME_LINUX)
	cd bin/linux_386 && tar -czf ../../dist/$(BINARY_NAME_LINUX)_$(VERSION)_linux_386.tar.gz $(BINARY_NAME_LINUX)
	cd bin/linux_arm64 && tar -czf ../../dist/$(BINARY_NAME_LINUX)_$(VERSION)_linux_arm64.tar.gz $(BINARY_NAME_LINUX)
	# Windows packages
	cd bin/windows_amd64 && zip -r ../../dist/$(BINARY_NAME_WINDOWS)_$(VERSION)_windows_amd64.zip $(BINARY_NAME_WINDOWS)
	cd bin/windows_386 && zip -r ../../dist/$(BINARY_NAME_WINDOWS)_$(VERSION)_windows_386.zip $(BINARY_NAME_WINDOWS)
	@echo "Distribution packages created in dist/ directory"

# Help information
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build       - Build for current platform"
	@echo "  clean       - Remove build artifacts"
	@echo "  deps        - Check and update dependencies"
	@echo "  test        - Run tests"
	@echo "  linux       - Build for Linux (amd64, 386, arm64)"
	@echo "  windows     - Build for Windows (amd64, 386)"
	@echo "  cross-build - Build for all supported platforms"
	@echo "  dist        - Create distribution packages"
	@echo "  help        - Show this help message"