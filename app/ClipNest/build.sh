#!/bin/bash
set -euo pipefail

cd "$(dirname "$0")"

echo "Building ClipNest..."
swift build -c release

APP_DIR="build/ClipNest.app/Contents"
mkdir -p "$APP_DIR/MacOS"

cp .build/release/ClipNest "$APP_DIR/MacOS/ClipNest"
cp Resources/Info.plist "$APP_DIR/Info.plist"

# Generate .icns from PNG
ICONSET="build/AppIcon.iconset"
mkdir -p "$ICONSET" "$APP_DIR/Resources"
sips -z 16 16     Resources/AppIcon.png --out "$ICONSET/icon_16x16.png"
sips -z 32 32     Resources/AppIcon.png --out "$ICONSET/icon_16x16@2x.png"
sips -z 32 32     Resources/AppIcon.png --out "$ICONSET/icon_32x32.png"
sips -z 64 64     Resources/AppIcon.png --out "$ICONSET/icon_32x32@2x.png"
sips -z 128 128   Resources/AppIcon.png --out "$ICONSET/icon_128x128.png"
sips -z 256 256   Resources/AppIcon.png --out "$ICONSET/icon_128x128@2x.png"
sips -z 256 256   Resources/AppIcon.png --out "$ICONSET/icon_256x256.png"
sips -z 512 512   Resources/AppIcon.png --out "$ICONSET/icon_256x256@2x.png"
sips -z 512 512   Resources/AppIcon.png --out "$ICONSET/icon_512x512.png"
sips -z 1024 1024 Resources/AppIcon.png --out "$ICONSET/icon_512x512@2x.png"
iconutil -c icns "$ICONSET" -o "$APP_DIR/Resources/AppIcon.icns"
rm -rf "$ICONSET"

codesign --force --deep -s - build/ClipNest.app

echo "Built: build/ClipNest.app"
