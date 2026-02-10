# ClipNest

A privacy-first clipboard manager for macOS and Linux. Keep track of your clipboard history, pin important clips, and never lose anything you've copied again.

## Features

- **Menu Bar App** - Native macOS SwiftUI app with Liquid Glass effects (macOS 26+) and material fallback (macOS 15+)
- **In-Memory Storage** - Recent clips (last 50) stored in RAM for instant access, no disk I/O
- **Pin Important Clips** - Mark clips to protect them from eviction
- **Real-Time Updates** - Unix socket IPC for instant synchronization between daemon and clients
- **CLI Interface** - Full command-line control over your clipboard history
- **Deduplication** - Automatically skips duplicate clips
- **Homebrew Support** - Install with `brew install` or `brew install --cask`

## Installation

### Homebrew (recommended)

```bash
brew tap AleCo3lho/clipnest
```

Install the GUI app (includes CLI + daemon):

```bash
brew install --cask clipnest-app
```

Or install just the CLI and daemon (no GUI):

```bash
brew install clipnest
brew services start clipnest
```

### Manual Build

Prerequisites: Go 1.23+, Swift 6.2+ (for the macOS app)

```bash
git clone https://github.com/AleCo3lho/clipnest.git
cd clipnest

# Build Go binaries
make build

# Build macOS menu bar app (requires macOS with Xcode 26)
make app

# Start the daemon
./bin/clipnestd
```

## Usage

### Menu Bar App

The ClipNest menu bar app lives in your macOS menu bar. It shows your clipboard history, lets you search, pin, and copy clips with a single click. The daemon starts automatically when you install via Homebrew.

### CLI Commands

| Command | Description |
|---------|-------------|
| `clipnest list [limit]` | List recent clips |
| `clipnest search <query>` | Search clips |
| `clipnest copy <id>` | Copy clip to clipboard |
| `clipnest pin <id>` | Pin clip |
| `clipnest unpin <id>` | Unpin clip |
| `clipnest pins` | List pinned clips |
| `clipnest clear` | Clear all clips |
| `clipnest version` | Show version |

### Quick Start (CLI)

```bash
# Start the daemon (or use brew services)
clipnestd &

# List recent clips
clipnest list

# Search clips
clipnest search "api"

# Copy a clip back to clipboard
clipnest copy 5

# Pin an important clip
clipnest pin 5
```

## Architecture

```
clipnest/
├── cmd/
│   ├── clipnest/              # CLI tool
│   └── clipnestd/             # Background daemon
├── internal/
│   ├── clipboard/             # Clipboard monitoring
│   ├── storage/               # In-memory LRU storage
│   ├── socket/                # Unix domain socket IPC
│   └── config/                # Configuration
├── app/ClipNest/              # macOS SwiftUI menu bar app
│   ├── Sources/
│   │   ├── Models/            # Clip data model
│   │   ├── Networking/        # Daemon socket client
│   │   ├── ViewModels/        # UI state management
│   │   ├── Views/             # SwiftUI views
│   │   └── Utilities/         # Helpers
│   └── Resources/             # Info.plist, app icon
└── homebrew/                  # Homebrew formula + cask templates
```

### How It Works

1. **clipnestd** (daemon) monitors the system clipboard, stores clips in an LRU cache (default: 50 clips), and serves them over a Unix socket at `/tmp/clipnest.sock`
2. **clipnest** (CLI) or **ClipNest.app** (menu bar) connects to the daemon socket to list, search, copy, and pin clips
3. Pinned clips are exempt from LRU eviction
4. All storage is in-memory only - nothing is written to disk

### Socket Protocol

Line-delimited JSON over Unix socket at `/tmp/clipnest.sock`:

```json
{"type":"new_clip","data":{"id":1,"content":"text","type":"text","timestamp":1234567890,"pinned":false}}
{"type":"copy_clip","data":{"id":1}}
{"type":"list","data":{"limit":100}}
{"type":"search","data":{"query":"api","limit":50}}
```

## Development

### Prerequisites

- Go 1.23+
- Swift 6.2+ / Xcode 26 (for macOS menu bar app)

### Build & Test

```bash
# Full quality gate (fmt, vet, lint, test)
make check

# Build everything
make build     # Go binaries
make app       # macOS menu bar app

# Run tests
make test

# See all targets
make help
```

### CI/CD

- **CI** (`ci.yml`): Runs on every push/PR - Go tests + linting (ubuntu), Swift build (macos-26)
- **Release** (`release.yml`): Triggers on `v*` tags - builds Go binaries + Swift app, creates GitHub release, updates Homebrew tap

## Future Plans

- [ ] Image clipboard support
- [ ] File path clipboard support
- [ ] Fuzzy search
- [ ] Global hotkey
- [ ] Export/import clips

## License

MIT License - feel free to use this for your own projects!
