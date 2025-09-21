# Makefile for Kick Clipper Go Version

# Build variables
BINARY_NAME=kick-clipper
BUILD_DIR=bin
MAIN_PATH=.

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build targets
.PHONY: all build clean test deps run help

all: clean deps build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME).exe $(MAIN_PATH)

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@if exist $(BUILD_DIR) rmdir /s $(BUILD_DIR)

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME).exe

run-no-proxy: build
	@echo "Running $(BINARY_NAME) in no-proxy mode..."
	./$(BUILD_DIR)/$(BINARY_NAME).exe --no-proxy

run-help: build
	@echo "Showing help..."
	./$(BUILD_DIR)/$(BINARY_NAME).exe --help

dev:
	@echo "Running in development mode..."
	$(GOCMD) run $(MAIN_PATH)

fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

help:
	@echo "Available targets:"
	@echo "  all       - Clean, download deps, and build"
	@echo "  build     - Build the binary"
	@echo "  clean     - Clean build artifacts"
	@echo "  deps      - Download dependencies"
	@echo "  test      - Run tests"
	@echo "  run       - Build and run"
	@echo "  run-no-proxy - Build and run in no-proxy mode"
	@echo "  run-help  - Build and show help"
	@echo "  dev       - Run in development mode"
	@echo "  fmt       - Format code"
	@echo "  vet       - Vet code"
	@echo "  help      - Show this help"