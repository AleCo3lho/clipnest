#!/bin/bash
set -euo pipefail

cd "$(dirname "$0")"

echo "Building ClipNest..."
swift build -c release

APP_DIR="build/ClipNest.app/Contents"
mkdir -p "$APP_DIR/MacOS"

cp .build/release/ClipNest "$APP_DIR/MacOS/ClipNest"
cp Resources/Info.plist "$APP_DIR/Info.plist"

echo "Built: build/ClipNest.app"
