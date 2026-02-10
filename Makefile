.PHONY: check fmt vet lint lint-install test build clean tidy app app-run app-clean dev help

# Default target: full quality gate
check: fmt vet lint test

## fmt: Format all Go files
fmt:
	go fmt ./...

## vet: Run Go static analysis
vet:
	go vet ./...

## lint: Run golangci-lint
lint:
	@which golangci-lint > /dev/null 2>&1 || { echo "Error: golangci-lint not installed. Run 'make lint-install' first."; exit 1; }
	golangci-lint run ./...

## lint-install: Install golangci-lint (pinned version)
lint-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$$(go env GOPATH)/bin" v1.62.2

## test: Run tests with race detection and coverage
test:
	go test -v -race -cover ./...

## build: Build binaries to bin/
build:
	@mkdir -p bin
	@if [ -d cmd/clipnest ]; then \
		echo "Building clipnest..."; \
		go build -o bin/clipnest ./cmd/clipnest; \
	else \
		echo "Skipping clipnest (cmd/clipnest not found)"; \
	fi
	@if [ -d cmd/clipnestd ]; then \
		echo "Building clipnestd..."; \
		go build -o bin/clipnestd ./cmd/clipnestd; \
	else \
		echo "Skipping clipnestd (cmd/clipnestd not found)"; \
	fi

## app: Build the macOS menu bar app
app:
	cd app/ClipNest && bash build.sh

## app-run: Build and launch the menu bar app
app-run: app
	open app/ClipNest/build/ClipNest.app

## app-clean: Remove Swift build artifacts
app-clean:
	rm -rf app/ClipNest/.build app/ClipNest/build

## dev: Build everything and run (daemon + app)
dev: build app
	@if [ -f bin/clipnestd ]; then \
		echo "Starting clipnestd..."; \
		./bin/clipnestd & \
	else \
		echo "Warning: bin/clipnestd not found, skipping daemon"; \
	fi
	@echo "Launching ClipNest.app..."
	@open app/ClipNest/build/ClipNest.app

## clean: Kill running processes and remove build artifacts
clean: app-clean
	@echo "Stopping running processes..."
	@pkill -x clipnestd || true
	@pkill -f ClipNest.app/Contents/MacOS/ClipNest || true
	rm -rf bin/
	go clean -testcache

## tidy: Clean up module dependencies
tidy:
	go mod tidy

## help: Show available targets
help:
	@echo "Available targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /' | sed 's/: /\t/'
