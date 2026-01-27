# ClipNest 🪹

Your cozy clipboard manager for macOS. Keep track of your clipboard history, pin important clips, and never lose anything you've copied again.

## Features

✅ **In-Memory Storage** - Recent clips (last 50) stored in RAM for instant access
✅ **Pin Important Clips** - Mark important clips for easy access
✅ **Real-Time Updates** - Unix socket for instant UI synchronization
✅ **CLI Interface** - Full command-line control
✅ **Privacy First** - Clips are only in memory (no disk I/O)
✅ **Deduplication** - Automatically skips duplicate clips

## Installation

## Installation

### Manual Build

#### Prerequisites
- Go 1.23+

#### macOS
```bash
# Clone repository
git clone https://github.com/AleCo3lho/clipnest.git
cd clipnest

# Build binaries
go build -o clipnest ./cmd/clipnest
go build -o clipnestd ./cmd/clipnestd

# Move to PATH
sudo mv clipnest clipnestd /usr/local/bin/

# Start daemon
clipnestd
```

#### Linux
```bash
# Clone repository
git clone https://github.com/AleCo3lho/clipnest.git
cd clipnest

# Build binaries
go build -o clipnest ./cmd/clipnest
go build -o clipnestd ./cmd/clipnestd

# Move to PATH
sudo mv clipnest clipnestd /usr/local/bin/

# Start daemon
clipnestd
```

### Build

```bash
# Build daemon and CLI
go build -o clipnestd ./cmd/clipnestd
go build -o clipnest ./cmd/clipnest
```

## Quick Start

### Build

```bash
# Build daemon and CLI
go build -o clipnestd ./cmd/clipnestd
go build -o clipnest ./cmd/clipnest

# Move to PATH
sudo cp clipnestd clipnest /usr/local/bin/
```

### Run

```bash
# Start the daemon
clipnestd

# In another terminal, list clips
clipnest list

# Search clips
clipnest search "api"

# Copy a clip
clipnest copy 5

# Pin an important clip
clipnest pin 5

# List pinned clips
clipnest pins
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `clipnest list [limit]` | List recent clips |
| `clipnest search <query>` | Search clips |
| `clipnest copy <id>` | Copy clip to clipboard |
| `clipnest pin <id>` | Pin clip (persist) |
| `clipnest unpin <id>` | Unpin clip |
| `clipnest pins` | List pinned clips |
| `clipnest clear` | Clear all clips |
| `clipnest version` | Show version |

## Architecture

```
clipnest/
├── cmd/
│   ├── clipnest/           ← CLI tool
│   └── clipnestd/          ← Background daemon
├── internal/
│   ├── clipboard/          ← Clipboard monitoring
│   ├── storage/           ← Hybrid storage (memory + SQLite)
│   ├── socket/            ← Unix domain socket
│   └── config/             ← Configuration
└── swift/                 ← Future Swift UI app
```

### Storage Strategy

- **In-Memory**: Recent clips (default: last 50)
  - Fast access
  - Auto-evicts oldest clips
  - Session-only by default
  - No disk I/O

- **SQLite**: Pinned clips only
  - User must explicitly pin
  - Persists across restarts
  - Searchable
  - Location: `~/Library/Application Support/ClipNest/clipnest.db`

### Socket Protocol

Unix socket at `/tmp/clipnest.sock` for real-time communication:

**Daemon → UI (events)**:
```json
{"type":"new_clip","data":{"id":1,"content":"text","type":"text","timestamp":1234567890,"pinned":false}}
{"type":"clipboard_changed","data":{"id":1}}
{"type":"clip_updated","data":{...}}
```

**UI → Daemon (commands)**:
```json
{"type":"copy_clip","data":{"id":1}}
{"type":"pin","data":{"id":1}}
{"type":"unpin","data":{"id":1}}
{"type":"list","data":{"limit":100}}
{"type":"search","data":{"query":"api","limit":50}}
```

## Development

### Prerequisites
- Go 1.23+
- SQLite3

### Setup
```bash
# Clone repository
git clone https://github.com/AleCo3lho/clipnest.git
cd clipnest

# Install dependencies
go mod tidy

# Run tests
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Build binaries
go build -o clipnest ./cmd/clipnest
go build -o clipnestd ./cmd/clipnestd
```

### Testing
```bash
# Run all tests
go test -v ./...

# Run tests with race detector
go test -race ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test -v ./internal/storage
```

### CI/CD
- **CI**: Runs on every push/PR - tests all packages
- **Release**: Triggers on version tags (v*) - builds macOS binaries

## Contributing

### Dependencies

```bash
go mod tidy
go mod download
```

### Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Future Plans

- [ ] Swift UI app for macOS
- [ ] Image clipboard support
- [ ] File path clipboard support
- [ ] Fuzzy search
- [ ] Global hotkey
- [ ] System tray icon
- [ ] Export/import clips
- [ ] Homebrew formula

## License

MIT License - feel free to use this for your own projects!

---

Built with ❤️ using Go
