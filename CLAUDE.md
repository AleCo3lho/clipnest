# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ClipNest is a privacy-first macOS/Linux clipboard manager written in Go. It features in-memory storage (no disk I/O by default), real-time synchronization via Unix domain sockets, and automatic deduplication. It consists of a CLI tool (`clipnest`) and a background daemon (`clipnestd`), though the `cmd/` entry points are not yet implemented.

## Build & Test Commands

A `Makefile` provides all development workflows. Run `make help` for available targets.

```bash
# Full quality gate (fmt → vet → lint → test)
make check

# Individual targets
make fmt          # Format all Go files
make vet          # Static analysis
make lint         # Run golangci-lint (install first with: make lint-install)
make test         # Tests with race detection and coverage
make build        # Build binaries to bin/
make clean        # Remove artifacts
make tidy         # go mod tidy

# Run tests for a specific package
go test -v ./internal/storage
```

CI uses `go test -v -race -cover ./...` on Go 1.23. Linter config is in `.golangci.yml`.

## Architecture

All application code lives under `internal/` in four packages:

- **storage** — Core data layer. `MemoryStore` implements an LRU cache using a doubly-linked list + hash map for O(1) operations. `Storage` wraps `MemoryStore` with capacity enforcement (default 50 clips), auto-eviction, and higher-level operations (Add, List, Pin, Search, etc.). Both use `sync.RWMutex` for thread safety.

- **clipboard** — `Monitor` polls the system clipboard at configurable intervals, detects changes by comparing content/type, and notifies via onChange callback. Uses `github.com/atotto/clipboard` for cross-platform access.

- **socket** — Unix domain socket server at `/tmp/clipnest.sock` for IPC between daemon and clients. Line-delimited JSON protocol. Broadcasts events to all connected clients. Thread-safe client management.

- **config** — Hardcoded defaults (50 clips max, socket path, data directory). Placeholder for future file-based configuration.

The central data model is `Clip` (ID, Content, Type, Timestamp, Pinned) defined in `storage/models.go`. Socket messages use `SocketMessage` with JSON serialization.

## Key Patterns

- Thread safety via `sync.RWMutex` throughout — no channels for coordination
- Observer pattern: clipboard monitor uses callbacks, not channels
- Tests use pure `testing.T` assertions (no test framework), with `setupTestStorage()` helpers creating a 5-clip-limit store
- Pinned clips are exempt from LRU eviction
