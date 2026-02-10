// swift-tools-version: 6.2

import PackageDescription

let package = Package(
    name: "ClipNest",
    platforms: [.macOS(.v26)],
    targets: [
        .executableTarget(
            name: "ClipNest",
            path: "Sources"
        ),
    ]
)
